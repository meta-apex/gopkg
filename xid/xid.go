package xid

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"os"
	"sort"
	"strconv"
	"sync/atomic"
	"time"
)

// Code inspired from mgo/bson ObjectId
// ID represents a unique request id
type ID [rawLen]byte

const (
	encodedLen = 20 // string encoded len
	rawLen     = 12 // binary raw len

	// encoding stores a custom version of the base32 encoding with lower case letters.
	encoding = "0123456789abcdefghijklmnopqrstuv"
)

var (
	// objectIDCounter is atomically incremented when generating a new ObjectId. It's
	// used as the counter part of an id. This id is initialized with a random value.
	objectIDCounter = randInt()

	// machineID is generated once and used in subsequent calls to the New* functions.
	machineID = readMachineID()

	// pid stores the current process id
	pid = os.Getpid()

	nilID ID

	// dec is the decoding map for base32 encoding
	dec [256]byte
)

func init() {
	for i := 0; i < len(dec); i++ {
		dec[i] = 0xFF
	}
	for i := 0; i < len(encoding); i++ {
		dec[encoding[i]] = byte(i)
	}

	// If /proc/self/cpuset exists and is not /, we can assume that we are in a
	// form of container and use the content of cpuset xor-ed with the PID in
	// order get a reasonable machine global unique PID.
	b, err := os.ReadFile("/proc/self/cpuset")
	if err == nil && len(b) > 1 {
		pid ^= int(crc32.ChecksumIEEE(b))
	}
}

// readMachineID generates a machine ID, derived from a platform-specific machine ID
// value, or else the machine's hostname, or else a randomly-generated number.
// It panics if all of these methods fail.
func readMachineID() []byte {
	// Allow env overrides for the machine id
	if id := readMachineIDFromEnv(); len(id) == 3 {
		return id
	}

	id := make([]byte, 3)
	hid, err := readPlatformMachineID()
	if err != nil || len(hid) == 0 {
		hid, err = os.Hostname()
	}
	if err == nil && len(hid) != 0 {
		hw := sha256.New()
		hw.Write([]byte(hid))
		copy(id, hw.Sum(nil))
	} else {
		// Fallback to rand number if machine id can't be gathered
		if _, randErr := rand.Reader.Read(id); randErr != nil {
			panic(fmt.Errorf("xid: cannot get hostname nor generate a random number: %v; %v", err, randErr))
		}
	}
	return id
}

func readMachineIDFromEnv() []byte {
	envMachineID := os.Getenv("XID_MACHINE_ID")
	if envMachineID == "" {
		return nil
	}

	num, err := strconv.Atoi(envMachineID)
	if err != nil {
		panic("XID_MACHINE_ID value is set to not a number")
	}

	if num < 0 || num > 0xFFFFFF {
		panic("XID_MACHINE_ID out of range for 3 bytes")
	}

	// Encode the number into big endian.
	return []byte{byte(num >> 16), byte(num >> 8), byte(num)}
}

// randInt generates a random uint32
func randInt() uint32 {
	b := make([]byte, 3)
	if _, err := rand.Reader.Read(b); err != nil {
		panic(fmt.Errorf("xid: cannot generate random number: %v;", err))
	}
	return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
}

// New generates a globally unique ID
func New() ID {
	return NewWithTime(time.Now())
}

// NewWithTime generates a globally unique ID with the passed in time
func NewWithTime(t time.Time) ID {
	var id ID
	// Timestamp, 4 bytes, big endian
	binary.BigEndian.PutUint32(id[:], uint32(t.Unix()))
	// Machine ID, 3 bytes
	id[4] = machineID[0]
	id[5] = machineID[1]
	id[6] = machineID[2]
	// Pid, 2 bytes, specs don't specify endianness, but we use big endian.
	id[7] = byte(pid >> 8)
	id[8] = byte(pid)
	// Increment, 3 bytes, big endian
	i := atomic.AddUint32(&objectIDCounter, 1)
	id[9] = byte(i >> 16)
	id[10] = byte(i >> 8)
	id[11] = byte(i)
	return id
}

// FromString reads an ID from its string representation
func FromString(id string) (ID, error) {
	i := &ID{}
	err := i.UnmarshalText([]byte(id))
	return *i, err
}

// String returns a base32 hex lowercased with no padding representation of the id (char set is 0-9, a-v).
func (id ID) String() string {
	text := make([]byte, encodedLen)
	encode(text, id[:])
	return string(text)
}

// Encode encodes the id using base32 encoding, writing 20 bytes to dst and return it.
func (id ID) Encode(dst []byte) []byte {
	encode(dst, id[:])
	return dst
}

// MarshalText implements encoding/text TextMarshaler interface
func (id ID) MarshalText() ([]byte, error) {
	text := make([]byte, encodedLen)
	encode(text, id[:])
	return text, nil
}

// MarshalJSON implements encoding/json Marshaler interface
func (id ID) MarshalJSON() ([]byte, error) {
	if id.IsNil() {
		return []byte("null"), nil
	}
	text := make([]byte, encodedLen+2)
	encode(text[1:encodedLen+1], id[:])
	text[0], text[encodedLen+1] = '"', '"'
	return text, nil
}

// encode by unrolling the stdlib base32 algorithm + removing all safe checks
func encode(dst, id []byte) {
	_ = dst[19]
	_ = id[11]

	dst[19] = encoding[(id[11]<<4)&0x1F]
	dst[18] = encoding[(id[11]>>1)&0x1F]
	dst[17] = encoding[(id[11]>>6)|(id[10]<<2)&0x1F]
	dst[16] = encoding[id[10]>>3]
	dst[15] = encoding[id[9]&0x1F]
	dst[14] = encoding[(id[9]>>5)|(id[8]<<3)&0x1F]
	dst[13] = encoding[(id[8]>>2)&0x1F]
	dst[12] = encoding[id[8]>>7|(id[7]<<1)&0x1F]
	dst[11] = encoding[(id[7]>>4)|(id[6]<<4)&0x1F]
	dst[10] = encoding[(id[6]>>1)&0x1F]
	dst[9] = encoding[(id[6]>>6)|(id[5]<<2)&0x1F]
	dst[8] = encoding[id[5]>>3]
	dst[7] = encoding[id[4]&0x1F]
	dst[6] = encoding[id[4]>>5|(id[3]<<3)&0x1F]
	dst[5] = encoding[(id[3]>>2)&0x1F]
	dst[4] = encoding[id[3]>>7|(id[2]<<1)&0x1F]
	dst[3] = encoding[(id[2]>>4)|(id[1]<<4)&0x1F]
	dst[2] = encoding[(id[1]>>1)&0x1F]
	dst[1] = encoding[(id[1]>>6)|(id[0]<<2)&0x1F]
	dst[0] = encoding[id[0]>>3]
}

// UnmarshalText implements encoding/text TextUnmarshaler interface
func (id *ID) UnmarshalText(text []byte) error {
	if len(text) != encodedLen {
		return ErrInvalidID
	}
	for _, c := range text {
		if dec[c] == 0xFF {
			return ErrInvalidID
		}
	}
	if !decode(id, text) {
		*id = nilID
		return ErrInvalidID
	}
	return nil
}

// UnmarshalJSON implements encoding/json Unmarshaler interface
func (id *ID) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "null" {
		*id = nilID
		return nil
	}
	// Check the slice length to prevent panic on passing it to UnmarshalText()
	if len(b) < 2 {
		return ErrInvalidID
	}
	return id.UnmarshalText(b[1 : len(b)-1])
}

// decode by unrolling the stdlib base32 algorithm + customized safe check.
func decode(id *ID, src []byte) bool {
	_ = src[19]
	_ = id[11]

	id[11] = dec[src[17]]<<6 | dec[src[18]]<<1 | dec[src[19]]>>4
	// check the last byte
	if encoding[(id[11]<<4)&0x1F] != src[19] {
		return false
	}
	id[10] = dec[src[16]]<<3 | dec[src[17]]>>2
	id[9] = dec[src[14]]<<5 | dec[src[15]]
	id[8] = dec[src[12]]<<7 | dec[src[13]]<<2 | dec[src[14]]>>3
	id[7] = dec[src[11]]<<4 | dec[src[12]]>>1
	id[6] = dec[src[9]]<<6 | dec[src[10]]<<1 | dec[src[11]]>>4
	id[5] = dec[src[8]]<<3 | dec[src[9]]>>2
	id[4] = dec[src[6]]<<5 | dec[src[7]]
	id[3] = dec[src[4]]<<7 | dec[src[5]]<<2 | dec[src[6]]>>3
	id[2] = dec[src[3]]<<4 | dec[src[4]]>>1
	id[1] = dec[src[1]]<<6 | dec[src[2]]<<1 | dec[src[3]]>>4
	id[0] = dec[src[0]]<<3 | dec[src[1]]>>2
	return true
}

// Time returns the timestamp part of the id.
// It's a runtime error to call this method with an invalid id.
func (id ID) Time() time.Time {
	// First 4 bytes of ObjectId is 32-bit big-endian seconds from epoch.
	secs := int64(binary.BigEndian.Uint32(id[0:4]))
	return time.Unix(secs, 0)
}

// Machine returns the 3-byte machine id part of the id.
// It's a runtime error to call this method with an invalid id.
func (id ID) Machine() []byte {
	return id[4:7]
}

// Pid returns the process id part of the id.
// It's a runtime error to call this method with an invalid id.
func (id ID) Pid() uint16 {
	return binary.BigEndian.Uint16(id[7:9])
}

// Counter returns the incrementing value part of the id.
// It's a runtime error to call this method with an invalid id.
func (id ID) Counter() int32 {
	b := id[9:12]
	// Counter is stored as big-endian 3-byte value
	return int32(uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2]))
}

// Value implements the driver.Valuer interface.
func (id ID) Value() (driver.Value, error) {
	if id.IsNil() {
		return nil, nil
	}
	b, err := id.MarshalText()
	return string(b), err
}

// Scan implements the sql.Scanner interface.
func (id *ID) Scan(value interface{}) (err error) {
	switch val := value.(type) {
	case string:
		return id.UnmarshalText([]byte(val))
	case []byte:
		return id.UnmarshalText(val)
	case nil:
		*id = nilID
		return nil
	default:
		return fmt.Errorf("xid: scanning unsupported type: %T", value)
	}
}

// IsNil Returns true if this is a "nil" ID
func (id ID) IsNil() bool {
	return id == nilID
}

// Alias of IsNil
func (id ID) IsZero() bool {
	return id.IsNil()
}

// NilID returns a zero value for `xid.ID`.
func NilID() ID {
	return nilID
}

// Bytes returns the byte array representation of `ID`
func (id ID) Bytes() []byte {
	return id[:]
}

// FromBytes convert the byte array representation of `ID` back to `ID`
func FromBytes(b []byte) (ID, error) {
	var id ID
	if len(b) != rawLen {
		return id, ErrInvalidID
	}
	copy(id[:], b)
	return id, nil
}

// Compare returns an integer comparing two IDs. It behaves just like `bytes.Compare`.
// The result will be 0 if two IDs are identical, -1 if current id is less than the other one,
// and 1 if current id is greater than the other.
func (id ID) Compare(other ID) int {
	return bytes.Compare(id[:], other[:])
}

type sorter []ID

func (s sorter) Len() int {
	return len(s)
}

func (s sorter) Less(i, j int) bool {
	return s[i].Compare(s[j]) < 0
}

func (s sorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Sort sorts an array of IDs inplace.
// It works by wrapping `[]ID` and use `sort.Sort`.
func Sort(ids []ID) {
	sort.Sort(sorter(ids))
}

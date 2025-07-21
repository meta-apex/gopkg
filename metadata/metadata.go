package metadata

import (
	"context"
	"strings"
)

type Metadata map[string]any

// New creates md from a given key-values map.
func New(mds ...map[string]any) Metadata {
	md := Metadata{}
	for _, m := range mds {
		for k, v := range m {
			md[k] = v
		}
	}
	return md
}

// Get returns the value associated with the passed key.
func (m Metadata) Get(key string) any {
	return m[strings.ToLower(key)]
}

// Set stores the key-value pair.
func (m Metadata) Set(key string, value any) {
	if key == "" || value == nil {
		return
	}

	m[strings.ToLower(key)] = value
}

func (m Metadata) IsExists(key string) bool {
	_, ok := m[key]
	return ok
}

// Range iterate over element in metadata.
func (m Metadata) Range(f func(k string, v any) bool) {
	for k, v := range m {
		if !f(k, v) {
			break
		}
	}
}

// Values returns a slice of values associated with the passed key.
func (m Metadata) Values(key string) any {
	return m[strings.ToLower(key)]
}

// Clone returns a deep copy of Metadata
func (m Metadata) Clone() Metadata {
	md := make(Metadata, len(m))
	for k, v := range m {
		md[k] = v
	}

	return md
}

func (m Metadata) FromStrMap(f map[string]string) {
	for k, v := range f {
		m[strings.ToLower(k)] = v
	}
}

type metadataContextKey struct{}

// NewContext creates a new context with metadata attached.
func NewContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, metadataContextKey{}, md)
}

// FromContext returns metadata from the given context.
func FromContext(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value(metadataContextKey{}).(Metadata)
	return md, ok
}

// MergeContext merges metadata to existing metadata, overwriting if specified.
func MergeContext(ctx context.Context, patchMd Metadata, overwrite bool) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	md, _ := ctx.Value(metadataContextKey{}).(Metadata)

	cmd := make(Metadata, len(md))
	for k, v := range md {
		cmd[k] = v
	}

	for k, v := range patchMd {
		if _, ok := cmd[k]; ok && !overwrite {
			// skip
		} else if v != "" {
			cmd[k] = v
		} else {
			delete(cmd, k)
		}
	}

	return context.WithValue(ctx, metadataContextKey{}, cmd)
}

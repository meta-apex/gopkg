package metadata

import (
	"context"
)

type Metadata map[string]any

// New creates md from a given key-values map.
func New(mds ...map[string]any) Metadata {
	md := make(map[string]any)
	for _, m := range mds {
		for k, v := range m {
			md[k] = v
		}
	}
	return md
}

// Get returns the value associated with the passed key.
func (m Metadata) Get(key string) (any, bool) {
	v, ok := m[key]
	return v, ok
}

// Set stores the key-value pair.
func (m Metadata) Set(key string, value any) {
	m[key] = value
}

// Clone returns a deep copy of Metadata
func (m Metadata) Clone() Metadata {
	md := make(Metadata, len(m))
	for k, v := range m {
		md[k] = v
	}
	return md
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

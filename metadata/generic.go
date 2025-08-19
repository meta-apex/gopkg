package metadata

import (
	"context"
)

// Generic represents metadata with generic key type
type Generic[K comparable] map[K]any

// NewGeneric creates generic metadata from given key-values maps
func NewGeneric[K comparable](mds ...map[K]any) Generic[K] {
	md := make(map[K]any)
	for _, m := range mds {
		for k, v := range m {
			md[k] = v
		}
	}
	return md
}

// Get returns the value associated with the passed key
func (m Generic[K]) Get(key K) (any, bool) {
	v, ok := m[key]
	return v, ok
}

// Set stores the key-value pair
func (m Generic[K]) Set(key K, value any) {
	m[key] = value
}

// Clone returns a deep copy of Generic
func (m Generic[K]) Clone() Generic[K] {
	md := make(Generic[K], len(m))
	for k, v := range m {
		md[k] = v
	}
	return md
}

// genericContextKey is used as context key for generic metadata
type genericContextKey[K comparable] struct{}

// NewGenericContext creates a new context with generic metadata attached
func NewGenericContext[K comparable](ctx context.Context, md Generic[K]) context.Context {
	return context.WithValue(ctx, genericContextKey[K]{}, md)
}

// FromGenericContext returns generic metadata from the given context
func FromGenericContext[K comparable](ctx context.Context) (Generic[K], bool) {
	md, ok := ctx.Value(genericContextKey[K]{}).(Generic[K])
	return md, ok
}

// MergeGenericContext merges generic metadata to existing metadata, overwriting if specified
func MergeGenericContext[K comparable](ctx context.Context, patchMd Generic[K], overwrite bool) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	md, _ := ctx.Value(genericContextKey[K]{}).(Generic[K])

	cmd := make(Generic[K], len(md))
	for k, v := range md {
		cmd[k] = v
	}

	for k, v := range patchMd {
		if _, ok := cmd[k]; ok && !overwrite {
			// skip
		} else if v != nil && v != "" {
			cmd[k] = v
		} else {
			delete(cmd, k)
		}
	}

	return context.WithValue(ctx, genericContextKey[K]{}, cmd)
}

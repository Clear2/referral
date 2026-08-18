// Package rediskey creates environment-scoped Redis keys.
package rediskey

import "strings"

const defaultPrefix = "prefix"

// RedisKey builds a namespaced Redis key.
type RedisKey interface {
	Key(parts ...string) string
}

type keyBuilder struct {
	environment string
	prefix      string
}

// New creates a Redis key builder scoped by environment and prefix.
func New(prefix, environment string) RedisKey {
	if prefix == "" {
		prefix = defaultPrefix
	}
	return &keyBuilder{environment: environment, prefix: prefix}
}

func (builder *keyBuilder) Key(parts ...string) string {
	all := make([]string, 0, len(parts)+2)
	all = append(all, builder.environment, builder.prefix)
	all = append(all, parts...)
	return strings.Join(all, ":")
}

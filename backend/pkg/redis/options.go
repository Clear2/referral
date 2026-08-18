package redis

// Option configures a Redis client.
type Option func(*Redis)

// MaxPoolSize sets the maximum number of socket connections.
func MaxPoolSize(size int) Option {
	return func(client *Redis) {
		client.maxPoolSize = size
	}
}

// Password sets the password used to authenticate with Redis. A non-empty
// value takes precedence over a password embedded in the connection URL.
func Password(password string) Option {
	return func(client *Redis) {
		client.password = password
	}
}

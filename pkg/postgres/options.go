package postgres

// Option -.
type Option func(*Postgres)

// MaxPoolSize -.
func MaxPoolSize(size int) Option {
	return func(c *Postgres) {
		c.maxPoolSize = size
	}
}

// AddLogger -.
func AddLogger(logger Logger) Option {
	return func(c *Postgres) {
		c.logger = logger
	}
}

package postgresql

import "time"

type Option func(*PostgreSQL)

func MaxPoolSize(size int) Option {
	return func(c *PostgreSQL) {
		c.maxPoolSize = size
	}
}

func ConnAttempts(attempts int) Option {
	return func(c *PostgreSQL) {
		c.connAttempts = attempts
	}
}

func ConnTimeout(timeout time.Duration) Option {
	return func(c *PostgreSQL) {
		c.connTimeout = timeout
	}
}

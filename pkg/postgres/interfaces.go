package postgres

type PGLogger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

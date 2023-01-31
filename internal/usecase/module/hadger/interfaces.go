package hadger

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=hadger_test

type Hadger interface{}

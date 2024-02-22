package core

type Interface interface {
	Init() error

	Load() error

	Close() error

	Display()
}

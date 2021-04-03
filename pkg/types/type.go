package types

type Migration interface {
	Up() error
	Down() error
}

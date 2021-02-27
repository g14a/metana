package _interface

type Migration interface {
	Up()
	Down()
}

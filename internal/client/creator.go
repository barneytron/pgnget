package client

import "os"

type Creatable interface {
	Create(name string) (*os.File, error)
}

type FileCreator struct{}

func NewCreator() Creatable {
	return &FileCreator{}
}

func (creator FileCreator) Create(name string) (*os.File, error) {
	return os.Create(name)
}

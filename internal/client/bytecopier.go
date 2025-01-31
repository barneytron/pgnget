package client

import "io"

type Copyable interface {
	Copy(dst io.Writer, src io.Reader) (written int64, err error)
}

type ByteCopier struct{}

func NewCopier() Copyable {
	return &ByteCopier{}
}

func (copier ByteCopier) Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	return io.Copy(dst, src)
}

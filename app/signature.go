package app

import (
	"hash/adler32"
	"io"
)

type signature struct {
}

func NewSignature(_ io.Reader, _ io.Writer) (*signature, error) {
	return &signature{}, nil
}

func hash() {
	a := adler32.New()
	a.Sum32()
}

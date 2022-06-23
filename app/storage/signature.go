package storage

import (
	"encoding/binary"
	"fmt"
	"io"
	"rollingHash/app"
)

type Signature struct {
	writer io.Writer
}

func NewSignatureStorage(out io.Writer) *Signature {
	return &Signature{writer: out}
}

func (s *Signature) Store(data *app.Signature) error {
	err := s.storeHeader(data)
	if err != nil {
		return fmt.Errorf("smth went wrong, %w", err)
	}

	err = s.storeData(data)
	if err != nil {
		return fmt.Errorf("smth went wrong, %w", err)
	}

	return nil
}

func (s *Signature) storeHeader(data *app.Signature) error {
	err := binary.Write(s.writer, binary.BigEndian, data.BlockSize)
	if err != nil {
		return fmt.Errorf("smth went wrong, %w", err)
	}

	err = binary.Write(s.writer, binary.BigEndian, data.RollingHashAlgo.Size())
	if err != nil {
		return fmt.Errorf("smth went wrong, %w", err)
	}

	err = binary.Write(s.writer, binary.BigEndian, data.StrongHashAlgo.Size())
	if err != nil {
		return fmt.Errorf("smth went wrong, %w", err)
	}
	return nil
}

func (s *Signature) storeData(data *app.Signature) error {
	for h, i := range data.Checksum() {
		err := binary.Write(s.writer, binary.BigEndian, h)
		if err != nil {
			return fmt.Errorf("smth went wrong, %w", err)
		}

		err = binary.Write(s.writer, binary.BigEndian, data.BlockStrongHash(i))
		if err != nil {
			return fmt.Errorf("smth went wrong, %w", err)
		}
	}
	return nil
}

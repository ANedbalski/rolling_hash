package storage

import (
	"encoding/binary"
	"fmt"
	"io"
	"rollingHash/app"
	"rollingHash/app/algo"
)

type Signature struct {
	writer io.Writer
	reader io.Reader
}

func NewSignatureStorage(w io.Writer, r io.Reader) *Signature {
	return &Signature{writer: w, reader: r}
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

func (s *Signature) Load() (*app.Signature, error) {
	blockSize, _, strongHashSize, err := s.loadHeader()
	if err != nil {
		return nil, err
	}
	sig := app.NewSignature(algo.NewAdler32(), algo.MD5, blockSize)
	err = s.loadData(sig, strongHashSize)
	return sig, err
}

func (s *Signature) storeHeader(data *app.Signature) error {
	err := binary.Write(s.writer, binary.BigEndian, data.BlockSize)
	if err != nil {
		return fmt.Errorf("smth went wrong, %w", err)
	}

	err = binary.Write(s.writer, binary.BigEndian, uint16(32))
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

func (s *Signature) loadHeader() (blockSize int64, rollHashSize uint16, strongHashSize uint16, err error) {
	err = binary.Read(s.reader, binary.BigEndian, &blockSize)
	if err != nil {
		return
	}
	err = binary.Read(s.reader, binary.BigEndian, &rollHashSize)
	if err != nil {
		return
	}
	err = binary.Read(s.reader, binary.BigEndian, &strongHashSize)
	if err != nil {
		return
	}
	return
}

func (s *Signature) loadData(sig *app.Signature, hashSize uint16) error {
	for {
		var weak uint32
		err := binary.Read(s.reader, binary.BigEndian, &weak)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		strong := make([]byte, hashSize)
		n, err := s.reader.Read(strong)
		if err != nil {
			return err
		}
		if n != int(hashSize) {
			return fmt.Errorf("unable to load hash data. read %d from %d", n, hashSize)
		}
		sig.AddHashRecord(weak, strong)
	}
	return nil
}

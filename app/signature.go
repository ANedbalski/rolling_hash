package app

import (
	md52 "crypto/md5"
	"encoding/binary"
	"fmt"
	"hash/adler32"
	"io"
)

const maxBlockSize = 512

type signature struct {
	blockSize   int
	rollingHash map[uint32]int64
	md5Hash     [][16]byte
}

func NewSignature(in io.Reader, out io.Writer) (*signature, error) {
	blockSize := maxBlockSize
	sig := &signature{
		blockSize:   blockSize,
		rollingHash: map[uint32]int64{},
		md5Hash:     make([][16]byte, 0),
	}
	for {
		block := make([]byte, blockSize)
		_, err := in.Read(block)
		if err == io.ErrUnexpectedEOF || err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("smth went wrong, %w", err)
		}

		checkSum := adler32.Checksum(block)
		md5 := md52.Sum(block)

		err = binary.Write(out, binary.BigEndian, checkSum)
		if err != nil {
			return nil, fmt.Errorf("smth went wrong, %w", err)
		}
		err = binary.Write(out, binary.BigEndian, md5)
		if err != nil {
			return nil, fmt.Errorf("smth went wrong, %w", err)
		}
		sig.rollingHash[checkSum] = int64(len(sig.md5Hash))
		sig.md5Hash = append(sig.md5Hash, md5)
	}

	return sig, nil
}

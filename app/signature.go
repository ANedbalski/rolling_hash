package app

import (
	"fmt"
	"io"
	"rollingHash/app/algo"
)

const maxBlockSize = 512

type Signature struct {
	RollingHashAlgo algo.RollingHash
	StrongHashAlgo  algo.StrongHash
	BlockSize       int64
	rollingHash     map[uint32]int64
	md5Hash         [][]byte
}

func NewSignature(r algo.RollingHash, s algo.StrongHash, blockSize int64) *Signature {
	return &Signature{
		RollingHashAlgo: r,
		StrongHashAlgo:  s,
		BlockSize:       blockSize,
		rollingHash:     make(map[uint32]int64),
		md5Hash:         make([][]byte, 0),
	}
}

func (s *Signature) Calc(in io.Reader) error {
	for {
		block := make([]byte, s.BlockSize)
		n, err := in.Read(block)
		if err == io.ErrUnexpectedEOF || err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("smth went wrong, %w", err)
		}
		checkSum := s.RollingHashAlgo.
			Reset().
			Write(block[:n]).
			Hash()

		md5 := s.StrongHashAlgo.Sum(block[:n])

		s.rollingHash[checkSum] = int64(len(s.md5Hash))
		s.md5Hash = append(s.md5Hash, md5)
	}

	return nil
}

func (s *Signature) BlockStrongHash(i int64) []byte {
	return s.md5Hash[i]
}

func (s *Signature) Find(sum uint32) (int64, bool) {
	i, ok := s.rollingHash[sum]
	return i, ok
}

func (s *Signature) Checksum() map[uint32]int64 {
	return s.rollingHash
}

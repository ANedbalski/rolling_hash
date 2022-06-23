package app

import (
	"bufio"
	"bytes"
	"github.com/balena-os/circbuf"
	"io"
	"rollingHash/app/algo"
)

type Delta struct {
	rollingHashAlgo algo.RollingHash
	strongHashAlgo  algo.StrongHash
}

func NewDelta(r algo.RollingHash, s algo.StrongHash) *Delta {
	return &Delta{
		rollingHashAlgo: r,
		strongHashAlgo:  s,
	}
}

func (d *Delta) Calc(sig *Signature, new io.Reader) (*DeltaData, error) {
	newBuf := bufio.NewReader(new)

	delta, err := NewDeltaData(sig.BlockSize)
	if err != nil {
		return nil, err
	}

	block, err := circbuf.NewBuffer(sig.BlockSize)
	if err != nil {
		return nil, err
	}

	weak := d.rollingHashAlgo.Reset()

	for {
		newByte, err := newBuf.ReadByte()
		if err == io.EOF {
			break
		}

		prevByte, err := blockWriteAndGetFirst(block, newByte)
		if err != nil {
			return nil, err
		}

		if block.TotalWritten() < sig.BlockSize {
			continue
		}

		if weak.Size() > 0 {
			weak.Roll(prevByte, newByte)
			err = delta.addNew(prevByte)
			if err != nil {
				return nil, err
			}
		} else {
			weak.Write(block.Bytes())
		}

		if ind, ok := d.validateHashes(sig, weak.Hash(), block.Bytes()); ok {
			delta.addCopy(uint64(ind*sig.BlockSize), uint64(sig.BlockSize))
			weak.Reset()
			block.Reset()
			continue
		}
	}

	weak.Reset().Write(block.Bytes())
	for i := int64(0); i < block.TotalWritten(); i++ {
		if ind, ok := d.validateHashes(sig, weak.Hash(), block.Bytes()[i:]); ok {
			delta.addCopy(uint64(ind*sig.BlockSize), uint64(sig.BlockSize))
			break
		}
		weak.Rollout(block.Bytes()[i])
		err = delta.addNew(block.Bytes()[i])
		if err != nil {
			return nil, err
		}
	}
	delta.flush()
	return delta, nil
}

func blockWriteAndGetFirst(block circbuf.Buffer, newByte byte) (prevByte byte, err error) {
	if block.TotalWritten() > 0 {
		prevByte, err = block.Get(0)
		if err != nil {
			return 0, err
		}
	}

	err = block.WriteByte(newByte)
	if err != nil {
		return 0, err
	}
	return prevByte, nil
}

func (d *Delta) validateHashes(sig *Signature, weak uint32, block []byte) (int64, bool) {
	ind, ok := sig.Find(weak)
	return ind,
		ok && bytes.Equal(sig.BlockStrongHash(ind), d.strongHashAlgo.Sum(block))
}

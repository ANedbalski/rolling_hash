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

	delta := NewDeltaData()

	block, err := circbuf.NewBuffer(sig.BlockSize)
	if err != nil {
		return nil, err
	}

	changes, err := circbuf.NewBuffer(sig.BlockSize)
	if err != nil {
		return nil, err
	}

	weak := d.rollingHashAlgo.Reset()
	var prevByte byte

	for {
		newByte, err := newBuf.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if block.TotalWritten() > 0 {
			prevByte, err = block.Get(0)
			if err != nil {
				return nil, err
			}
		}

		err = block.WriteByte(newByte)
		if err != nil {
			return nil, err
		}
		if block.TotalWritten() < sig.BlockSize {
			continue
		}

		if weak.Size() > 0 {
			weak.Roll(prevByte, newByte)

			err = changes.WriteByte(prevByte)
			if err != nil {
				return nil, err
			}
			if changes.TotalWritten() >= sig.BlockSize {
				delta.addNew(changes.Bytes())
				changes.Reset()
			}
		} else {
			weak.Write(block.Bytes())
		}

		if ind, ok := d.validateHashes(sig, weak.Hash(), block.Bytes()); ok {
			if changes.TotalWritten() > 0 {
				delta.addNew(changes.Bytes())
				changes.Reset()
			}
			delta.addCopy(uint64(ind*sig.BlockSize), uint64(sig.BlockSize))
			weak.Reset()
			block.Reset()
			continue
		}
	}

	if block.TotalWritten() > 0 {
		weak.Reset().Write(block.Bytes())
		for i := int64(0); i < sig.BlockSize; i++ {
			if ind, ok := d.validateHashes(sig, weak.Hash(), block.Bytes()[i:]); ok {
				if changes.TotalWritten() > 0 {
					delta.addNew(changes.Bytes())
					changes.Reset()
				}
				delta.addCopy(uint64(ind*sig.BlockSize), uint64(sig.BlockSize))
				break
			}
			prevByte = block.Bytes()[i]
			if err != nil {
				return nil, err
			}
			if err != nil {
				return nil, err
			}
			weak.Rollout(prevByte)
			err = changes.WriteByte(prevByte)
			if err != nil {
				return nil, err
			}
			if changes.TotalWritten() >= sig.BlockSize {
				delta.addNew(changes.Bytes())
				changes.Reset()
			}
		}
	}
	if changes.TotalWritten() > 0 {
		delta.addNew(changes.Bytes())
		changes.Reset()
	}
	return delta, nil
}

func (d *Delta) validateHashes(sig *Signature, weak uint32, block []byte) (int64, bool) {
	ind, ok := sig.Find(weak)
	return ind,
		ok && bytes.Equal(sig.BlockStrongHash(ind), d.strongHashAlgo.Sum(block))
}

func (d *Delta) readBlock() error {
	return nil
}

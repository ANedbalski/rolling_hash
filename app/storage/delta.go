package storage

import (
	"fmt"
	"io"
	"rollingHash/app"
)

type Delta struct {
	writer io.Writer
}

func NewDeltaStorage(out io.Writer) *Delta {
	return &Delta{writer: out}
}

func (d *Delta) Store(rs []app.Record) error {
	for _, r := range rs {
		err := d.storeFuncFactory(r.Op)(r)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Delta) storeAddOp(r app.Record) error {
	_, err := d.writer.Write([]byte(fmt.Sprintf("+%s\n", r.Data)))
	return err
}

func (d *Delta) storeCopyOp(r app.Record) error {
	_, err := d.writer.Write([]byte(fmt.Sprintf("=%d%d\n", r.Pos, r.Len)))
	return err
}

func (d *Delta) storeFuncFactory(op app.Op) func(app.Record) error {
	switch op {
	case app.OP_COPY:
		return d.storeCopyOp
	default:
		return d.storeAddOp
	}
}

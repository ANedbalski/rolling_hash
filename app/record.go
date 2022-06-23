package app

import (
	"github.com/balena-os/circbuf"
)

type Op int

const (
	OP_ADD  = 0
	OP_COPY = 1
)

type Record struct {
	Op   Op
	Pos  uint64
	Len  uint64
	Data []byte
}

type DeltaData struct {
	records   []Record
	blockSize int64
	newData   circbuf.Buffer
}

func NewDeltaData(s int64) (*DeltaData, error) {
	b, err := circbuf.NewBuffer(s)
	if err != nil {
		return nil, err
	}
	return &DeltaData{
		records:   []Record{},
		blockSize: s,
		newData:   b,
	}, nil
}

func (d *DeltaData) addCopy(pos, len uint64) {
	d.flush()
	d.records = append(d.records, Record{Op: OP_COPY, Pos: pos, Len: len})
}

func (d *DeltaData) addNew(b byte) error {
	err := d.newData.WriteByte(b)
	if err != nil {
		return err
	}
	if d.newData.TotalWritten() >= d.blockSize {
		d.flush()
	}
	return nil
}

func (d *DeltaData) flush() {
	if d.newData.TotalWritten() == 0 {
		return
	}
	b := make([]byte, d.newData.TotalWritten())
	copy(b, d.newData.Bytes())
	d.records = append(d.records, Record{Op: OP_ADD, Data: b})
	d.newData.Reset()
}

func (d *DeltaData) GetRecords() []Record {
	return d.records
}

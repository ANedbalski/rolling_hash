package app

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
	records []Record
}

func NewDeltaData() *DeltaData {
	return &DeltaData{records: []Record{}}
}

func (d *DeltaData) addCopy(pos, len uint64) {
	d.records = append(d.records, Record{Op: OP_COPY, Pos: pos, Len: len})
}

func (d *DeltaData) addNew(block []byte) {
	b := make([]byte, len(block))
	copy(b, block)
	d.records = append(d.records, Record{Op: OP_ADD, Data: b})
}

func (d *DeltaData) getRecords() []Record {
	return d.records
}

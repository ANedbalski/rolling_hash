package app

import (
	"github.com/stretchr/testify/assert"
	"rollingHash/app/algo"
	"strings"
	"testing"
)

func TestDelta_Calc(t *testing.T) {
	testCases := []struct {
		name      string
		old       string
		new       string
		blockSize int64
		sig       *Signature
		delta     []Record
	}{
		{
			name:      "versions are the same",
			old:       "AAAA AAAA AAAA AAAA",
			new:       "AAAA AAAA AAAA AAAA",
			blockSize: 10,
			delta: []Record{
				{Op: OP_COPY, Pos: 0, Len: 10},
				{Op: OP_COPY, Pos: 10, Len: 10},
			},
		},
		{
			name:      "versions are the same 2 block with full size",
			old:       "AAAA AAAA AAAA AAAAB",
			new:       "AAAA AAAA AAAA AAAAB",
			blockSize: 10,
			delta: []Record{
				{Op: OP_COPY, Pos: 0, Len: 10},
				{Op: OP_COPY, Pos: 10, Len: 10},
			},
		},
		{
			name:      "Insert data in the beginning less than 1 block",
			old:       "AAAA AAAA AAAA AAAA",
			new:       "BBB AAAA AAAA AAAA AAAA",
			blockSize: 10,
			delta: []Record{
				{Op: OP_ADD, Data: []byte("BBB ")},
				{Op: OP_COPY, Pos: 0, Len: 10},
				{Op: OP_COPY, Pos: 10, Len: 10},
			},
		},
		{
			name:      "Insert data in the beginning longer than 1 block",
			old:       "ABCDEFG",
			new:       "123456ABCDEFG",
			blockSize: 4,
			delta: []Record{
				{Op: OP_ADD, Data: []byte("1234")},
				{Op: OP_ADD, Data: []byte("56")},
				{Op: OP_COPY, Pos: 0, Len: 4},
				{Op: OP_COPY, Pos: 4, Len: 4},
			},
		},
		{
			name:      "Changes inside block",
			old:       "ABCDEFG",
			new:       "A12DEFG",
			blockSize: 4,
			delta: []Record{
				{Op: OP_ADD, Data: []byte("A12D")},
				{Op: OP_COPY, Pos: 4, Len: 4},
			},
		},
		{
			name:      "New data at the end",
			old:       "ABCDEFG",
			new:       "ABCDEFG5",
			blockSize: 4,
			delta: []Record{
				{Op: OP_COPY, Pos: 0, Len: 4},
				{Op: OP_ADD, Data: []byte("EFG5")},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			//Create signature for the old string
			sig := NewSignature(algo.NewAdler32(), algo.MD5, tt.blockSize)
			err := sig.Calc(strings.NewReader(tt.old))
			if !assert.Nil(t, err) {
				return
			}

			// produce and assert delta
			delta, err := NewDelta(algo.NewAdler32(), algo.MD5).Calc(sig, strings.NewReader(tt.new))
			if !assert.Nil(t, err) {
				return
			}
			assert.Equal(t, tt.delta, delta.GetRecords())
		})
	}

}

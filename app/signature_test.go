package app

import (
	"bytes"
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"rollingHash/app/algo"
	"testing"
)

func TestSignature_Calc(t *testing.T) {
	testCases := []struct {
		name      string
		data      []byte
		blockSize int64
		rolling   map[uint32]int64
		strong    [][]byte
	}{
		{
			name:      "1 block",
			data:      []byte("ABCDEFG"),
			blockSize: 8,
			rolling: map[uint32]int64{
				0x075b01dd: 0,
			},
			strong: [][]byte{
				[]byte("bb747b3df3130fe1ca4afa93fb7d97c9"),
			},
		},
		{
			name:      "2 full blocks",
			data:      []byte("1234567 12345678"),
			blockSize: 8,
			rolling: map[uint32]int64{
				0x0728018d: 0,
				0x074001a5: 1,
			},
			strong: [][]byte{
				[]byte("6a59d1ee75e8c6d7ad8b74b7759799b8"),
				[]byte("25d55ad283aa400af464c76d713c07ad"),
			},
		},
		{
			name:      "2 blocks. full and small",
			data:      []byte("1234567 12"),
			blockSize: 8,
			rolling: map[uint32]int64{
				0x0728018d: 0,
				0x00960064: 1,
			},
			strong: [][]byte{
				[]byte("6a59d1ee75e8c6d7ad8b74b7759799b8"),
				[]byte("c20ad4d76fe97759aa27a0c99bff6710"),
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			sig := NewSignature(algo.NewAdler32(), algo.MD5, tt.blockSize)
			err := sig.Calc(bytes.NewReader(tt.data))
			if !assert.Nil(t, err) {
				return
			}
			md5 := make([][]byte, len(tt.strong))
			for i, h := range tt.strong {
				md5[i], _ = hex.DecodeString(string(h))
			}
			assert.Equal(t, tt.rolling, sig.rollingHash)
			assert.Equal(t, md5, sig.md5Hash)
		})
	}
}

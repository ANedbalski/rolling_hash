package algo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdler32_Write(t *testing.T) {
	testCases := []struct {
		data []byte
		hash uint32
	}{
		{
			data: []byte("test string"),
			hash: 0x1ac00478,
		},
		{
			data: []byte("aaaaaaaaa"),
			hash: 0x1116036a,
		},
		{
			data: []byte("The real-world use case for this type of construct could be a distributed file storage system. This reduces the need for bandwidth and storage"),
			hash: 0x8124342d,
		},
		{
			data: []byte("12"),
			hash: 0x00960064,
		},
	}

	for _, tt := range testCases {
		t.Run(string(tt.data), func(t *testing.T) {
			assert.Equal(t, tt.hash, NewAdler32().Write(tt.data).Hash())
		})
	}
}

func TestAdler32_Roll(t *testing.T) {
	testCases := []struct {
		startingData []byte
		out          byte
		in           byte
		hashBefore   uint32
		hashAfter    uint32
	}{
		{
			startingData: []byte("12345678"),
			out:          []byte("1")[0],
			in:           []byte("9")[0],
			hashBefore:   0x074001a5,
			hashAfter:    0x076401ad,
		},
		{
			startingData: []byte("98765432"),
			out:          []byte("9")[0],
			in:           []byte("1")[0],
			hashBefore:   0x07b801ad,
			hashAfter:    0x079401a5,
		},
	}

	for _, tt := range testCases {
		t.Run(string(tt.startingData), func(t *testing.T) {
			a := NewAdler32().Write(tt.startingData)

			assert.Equal(t, tt.hashBefore, a.Hash())

			_ = a.Roll(tt.out, tt.in)

			assert.Equal(t, tt.hashAfter, a.Hash())
		})
	}
}

func TestAdler32_Rollout(t *testing.T) {
	testCases := []struct {
		name         string
		startingData []byte
		out          byte
		hashBefore   uint32
		hashAfter    uint32
	}{
		{
			name:         "remove 1st",
			startingData: []byte("12345678"),
			out:          []byte("1")[0],
			hashBefore:   0x074001a5,
			hashAfter:    0x05b70174,
		},
		{
			name:         "long data so possible negative numbers in calculations",
			startingData: []byte("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~"),
			out:          []byte("~")[0],
			hashBefore:   0x0102ec41,
			hashAfter:    0x14b2ebc3,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAdler32().Write(tt.startingData)

			assert.Equal(t, tt.hashBefore, a.Hash())

			_ = a.Rollout(tt.out)

			assert.Equal(t, tt.hashAfter, a.Hash())
		})
	}
}

func TestAdler32_Rollin(t *testing.T) {
	testCases := []struct {
		name         string
		startingData []byte
		in           byte
		hashBefore   uint32
		hashAfter    uint32
	}{
		{
			name:         "add 1st",
			startingData: []byte(""),
			in:           []byte("A")[0],
			hashBefore:   0x1,
			hashAfter:    0x00420042,
		},
		{
			name:         "add next",
			startingData: []byte("Any data"),
			in:           []byte("A")[0],
			hashBefore:   0x0c8402e3,
			hashAfter:    0x0fa80324,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAdler32().Write(tt.startingData)

			assert.Equal(t, tt.hashBefore, a.Hash())

			_ = a.Rollin(tt.in)

			assert.Equal(t, tt.hashAfter, a.Hash())

		})
	}

}

func TestAdler32_Reset(t *testing.T) {
	a := NewAdler32()
	_ = a.Write([]byte("any data"))
	assert.NotZero(t, a.Hash())

	a.Reset()

	assert.Equal(t, uint32(1), a.Hash())
}

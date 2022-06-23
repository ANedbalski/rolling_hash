package algo

type RollingHash interface {
	Write([]byte) RollingHash
	Roll(out, in byte) RollingHash
	Rollin(in byte) RollingHash
	Rollout(out byte) RollingHash
	Hash() uint32
	Reset() RollingHash
	Size() uint64
}

type StrongHash interface {
	Sum([]byte) []byte
	Size() uint16
}

type StrongHashImpl struct {
	f    func([]byte) []byte
	size uint16
	name string
}

func (s StrongHashImpl) Sum(b []byte) []byte {
	return s.f(b)
}

func (s StrongHashImpl) Size() uint16 {
	return s.size
}

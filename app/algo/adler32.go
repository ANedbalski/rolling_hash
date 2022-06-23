package algo

const Adler32Name = "adler32"

const (
	BASE_MOD = 65521
)

type Adler32 struct {
	cnt    uint64
	s1, s2 uint32
}

func NewAdler32() RollingHash {
	return &Adler32{}
}

func (a *Adler32) Write(b []byte) RollingHash {
	a.s1 = 1
	for _, i := range b {
		a.cnt++
		a.s1 = (a.s1 + uint32(i)) % BASE_MOD
		a.s2 = (a.s2 + a.s1) % BASE_MOD
	}
	return a
}

func (a *Adler32) Roll(out, in byte) RollingHash {
	a.s1 += uint32(in) - uint32(out)
	a.s1 %= BASE_MOD
	a.s2 = (a.s2 - (uint32(a.cnt) * uint32(out)) + a.s1 - 1) % BASE_MOD
	return a
}

func (a *Adler32) Rollout(out byte) RollingHash {
	a.s1 -= uint32(out)
	a.s1 %= BASE_MOD
	a.s2 = (a.s2 + BASE_MOD - uint32(a.cnt)*uint32(out) - 1) % BASE_MOD
	a.cnt--
	return a
}

func (a *Adler32) Hash() uint32 {
	return (a.s2 << 16) | a.s1
}

func (a *Adler32) Reset() RollingHash {
	a.cnt = 0
	a.s1 = 0
	a.s2 = 0
	return a
}

func (a *Adler32) Size() uint64 {
	return a.cnt
}

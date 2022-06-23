package algo

import "crypto/md5"

var MD5 = &StrongHashImpl{
	f: func(in []byte) []byte {
		sum := md5.Sum(in)
		return sum[:]
	},
	size: 16,
}

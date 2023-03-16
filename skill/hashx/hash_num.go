package hashx

import "github.com/spaolacci/murmur3"

func Hash(data []byte) uint64 {
	return murmur3.Sum64(data)
}

func Sum64(data []byte) uint64 {
	return murmur3.Sum64(data)
}

func Sum32(data []byte) uint32 {
	return murmur3.Sum32(data)
}

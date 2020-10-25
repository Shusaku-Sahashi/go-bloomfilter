package src

import (
	"math"
)

type Bloom struct {
	dataset []uint64
	numHash uint64
	size uint64
	shift uint64
}

func NewBloom(entries, wr float64) *Bloom {
	var filterSize, numHash uint64
	if wr < 1 {
		filterSize, numHash = calcParams(entries, wr)
	} else {
		filterSize, numHash = uint64(entries), uint64(wr)
	}

	size, exponent := getSize(filterSize)
	return &Bloom{
		dataset: make([]uint64, size>>6),
		numHash: numHash,
		shift: 64 - exponent,
		size: size - 1,
	}
}

func getSize(entries uint64) (uint64, uint64) {
	if entries < uint64(514) {
		entries = uint64(514)
	}

	var exponent uint64
	for i := uint64(1); i < entries; i <<= 1 {
		exponent++
	}

	return entries, exponent
}

func calcParams(numEntry, wr float64) (uint64, uint64) {
	filterSize := -1 * numEntry * math.Log(wr) / math.Pow(math.Log(2), 2)
	numHash := math.Ceil(filterSize * math.Log(2) / numEntry)
	return uint64(filterSize), uint64(numHash)
}

func (b *Bloom) nHash(h, l, i, m uint64) uint64 {
	return (l + i * h) % m
}

func (b *Bloom) Add(hash uint64) {
	h := hash >> 6
	l := hash << 6 >> 6
	for i := uint64(0); i < b.numHash; i++ {
		b.set(b.nHash(h, l, i, b.size))
	}
}

func (b *Bloom) Exist(hash uint64) bool {
	h := hash >> 6
	l := hash << 6 >> 6
	for i := uint64(0); i < b.numHash; i++ {
		if b.isSet(b.nHash(h, l, i, b.size)) == false {
			return false
		}
	}

	return true
}

func (b *Bloom) set(idx uint64) {
	b.dataset[idx>>6] |= 1<<(idx%64)
}

func (b *Bloom) isSet(idx uint64) bool {
	f := (b.dataset[idx>>6] >> (idx%64)) & 1
	return f == 1
}


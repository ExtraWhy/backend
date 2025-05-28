package bitvector

var ThisFile = "bitvector.go"

type Bitvector struct {
	bset  []uint64
	bsize int
}

func (bv *Bitvector) NewBitvector(sz int) {
	bv.bset = make([]uint64, sz+1)
	bv.bsize = (sz + 1)
}

func (bv *Bitvector) Add(item int) {
	bv.bset[item>>6] |= (1 << (item % 64))
}

func (bv *Bitvector) Contains(item int) bool {
	return (bv.bset[(item>>6)] & (1 << (item % 64))) != 0
}

func (bv *Bitvector) Indices() []uint8 {
	var a []uint8
	for i := 0; i < bv.bsize*64; i++ {
		if bv.Contains(i) {
			a = append(a, uint8(i))
		}
	}
	return a
}

func (bv *Bitvector) Reset() {
	for i := 0; i < len(bv.bset); i++ {
		bv.bset[i] = 0
	}
}

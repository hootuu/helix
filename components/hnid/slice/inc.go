package slice

type AutoIncSeq interface {
	Next() uint64
}

type AutoInc struct {
	length   uint8
	sequence AutoIncSeq
}

func NewAutoInc(length uint8, sequence AutoIncSeq) *AutoInc {
	return &AutoInc{
		length:   length,
		sequence: sequence,
	}
}

func (ais *AutoInc) Build() (uint64, uint8) {
	value := ais.sequence.Next()
	length := ais.length
	return value, length
}

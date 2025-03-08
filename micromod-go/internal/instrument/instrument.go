package instrument

const (
	FP_SHIFT = 15
	FP_ONE   = 1 << FP_SHIFT
	FP_MASK  = FP_ONE - 1
)

type Instrument struct {
	name                                    string
	volume, fineTune, loopStart, loopLength int
	sampleData                              []byte // todo: new byte[ 0 ] ???
}

func (i *Instrument) Audio(
	sampleIdx, sampleFrac, step, leftGain, rightGain int,
	mixBuf []int,
	offset, count int, interpolation bool,
) {
	// todo
}

func (i *Instrument) NormalizeSampleIdx(sampleIdx int) int {
	loopOffset := sampleIdx - i.loopStart
	if loopOffset > 0 {
		sampleIdx = i.loopStart
		if i.loopLength > 1 {
			sampleIdx += loopOffset % i.loopLength
		}
	}
	return sampleIdx
}

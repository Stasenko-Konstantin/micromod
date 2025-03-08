package module

import (
	"micromod/internal/instrument"
	"micromod/internal/pattern"
)

const (
	C2_PAL  = 8287
	C2_NTSC = 8363
)

type Module struct {
	modFilepath string

	songName                                 string
	sequenceLength, restartPos, c2Rate, gain int
	sequence                                 []byte
	patterns                                 []*pattern.Pattern
	instruments                              []*instrument.Instrument
}

func New(modFilepath string) *Module {
	return &Module{modFilepath: modFilepath}
}

func (m *Module) NumChannels() int {
	return m.patterns[0].NumChannels()
}

// todo
func (m *Module) ModuleInfo() string {
	return ""
}

// todo
func (m *Module) SequenceLength() int {
	return 0
}

// todo
func (m *Module) NumChannel() int {
	return 0
}

// todo
func (m *Module) SequenceEntry(seqIdx int) int {
	return 0
}

// todo
func (m *Module) Pattern(patIdx int) *pattern.Pattern {
	return nil
}

// todo
func (m *Module) Instrument(i int) *instrument.Instrument {
	return nil
}

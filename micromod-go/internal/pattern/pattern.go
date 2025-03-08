package pattern

import "micromod/internal/note"

const (
	NumRows = 64
)

type Pattern struct {
	numChannels int
	patternData []byte
}

func (p *Pattern) NumChannels() int {
	return p.numChannels
}

// todo
func (p *Pattern) Note(row, ch int, note *note.Note) *note.Note {
	return nil
}

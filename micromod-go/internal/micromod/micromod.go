/*
micromod - done
todo     - other functionality
*/
package micromod

import (
	"fmt"
	"micromod/internal/channel"
	"micromod/internal/module"
	"micromod/internal/note"
	"micromod/internal/player"
	"micromod/pattern"
)

type Micromod struct {
	module                                    *module.Module
	player                                    *player.Player
	rampBuf                                   []int
	note                                      *note.Note
	channels                                  []*channel.Channel
	sampleRate                                int
	seqPos, breakSeqPos, rowV, nextRow, tickV int
	speed, tempo, plCount, plChannel          int
	interpolation                             bool
	playCount                                 [][]byte
}

func New(modFilepath string, interpolation, loop bool) *Micromod {
	module := module.New(modFilepath)
	mm := &Micromod{
		module:    module,
		player:    player.New(interpolation, loop),
		rampBuf:   make([]int, 128),
		note:      note.New(),
		playCount: make([][]byte, module.SequenceLength()),
		channels:  make([]*channel.Channel, module.NumChannel()),
	}
	mm.setSequencePos(0)
	return mm
}

// todo
func (m *Micromod) Run(clCh chan struct{}) {
	clCh <- struct{}{}
}

func (m *Micromod) ModuleInfo() string {
	return m.module.ModuleInfo()
}

// SampleRate - return the sampling rate of playback
func (m *Micromod) SampleRate() int {
	return m.sampleRate
}

// SetSampleRate - set the sampling rate of playback
func (m *Micromod) SetSampleRate(rate int) error {
	// Use with Module.c2Rate to adjust the tempo of playback
	// To play at half speed, multiply both the samplingRate and Module.c2Rate by 2
	if rate < 8000 || rate > 128000 {
		return fmt.Errorf("rate must be between 8000 and 128000")
	}
	m.sampleRate = rate
	return nil
}

// SetInterpolation - Enable or disable the linear interpolation filter
func (m *Micromod) SetInterpolation(interpolation bool) {
	m.interpolation = interpolation
}

// MixBufferLength - Return the length of the buffer required by GetAudio()
func (m *Micromod) MixBufferLength() int {
	return (m.calculateTickLen(32, 128000) + 65) * 4
}

// Row - Get the current row position
func (m *Micromod) Row() int {
	return m.rowV
}

// SequencePos - Get the current pattern position in the sequence
func (m *Micromod) SequencePos() int {
	return m.seqPos
}

// setSequencePos - Set the pattern in the sequence to play. The tempo is reset to the default
func (m *Micromod) setSequencePos(pos int) {
	if pos >= m.module.SequenceLength() {
		pos = 0
	}
	m.breakSeqPos = pos
	m.nextRow = 0
	m.tickV = 1
	m.speed = 6
	m.tempo = 125
	m.plCount = -1
	m.plChannel = -1
	for i := 0; i < len(m.playCount); i++ {
		m.playCount[i] = make([]byte, pattern.NumRows)
	}
	for i := 0; i < len(m.channels); i++ {
		m.channels[i] = channel.New(m.module, i)
	}
	m.tick()
}

// CalculateSongDuration - Returns the song duration in samples at the current sampling rate
func (m *Micromod) CalculateSongDuration() int {
	duration := 0
	m.setSequencePos(0)
	songEnd := false
	for !songEnd {
		duration += m.calculateTickLen(m.tempo, m.sampleRate)
		songEnd = m.tick()
	}
	m.setSequencePos(0)
	return duration
}

// Audio - Generate audio
// The number of samples placed into outputBuf is returned
// The output buffer length must be at least that returned by getMixBufferLength()
// A "sample" is a pair of 16-bit integer amplitudes, one for each of the stereo channels
func (m *Micromod) Audio(outputBuf []int) int {
	tickLen := m.calculateTickLen(m.tempo, m.sampleRate)

	// Clear output buffer
	for i, end := 0, (tickLen+65)*4; i < end; i++ {
		outputBuf[i] = 0
	}

	// Resample
	for _, ch := range m.channels {
		ch.Resample(outputBuf, 0, (tickLen+65)*2, m.sampleRate*2, m.interpolation)
		ch.UpdateSampleIdx(tickLen*2, m.sampleRate*2)
	}
	m.downsample(outputBuf, tickLen+64)
	m.volumeRamp(outputBuf, tickLen)
	m.tick()
	return tickLen
}

func (m *Micromod) calculateTickLen(tempo, samplingRate int) int {
	return (samplingRate * 5) / (tempo * 2)
}

func (m *Micromod) volumeRamp(mixBuf []int, tickLen int) {
	rampRate := 256 * 2048 / m.sampleRate
	for i, a1 := 0, 0; a1 < 256; i, a1 = i+2, a1+rampRate {
		a2 := 256 - a1
		mixBuf[i] = (mixBuf[i]*a1 + m.rampBuf[i]*a2) >> 8
		mixBuf[i+1] = (mixBuf[i+1]*a1 + m.rampBuf[i+1]*a2) >> 8
	}
	copy(m.rampBuf, mixBuf[tickLen*2:tickLen*2+128]) // todo: maybe panic
}

func (m *Micromod) downsample(buf []int, count int) {
	// 2:1 downsampling with simple but effective anti-aliasing. Buf must contain count * 2 + 1 stereo samples
	outLen := count * 2
	for inIdx, outIdx := 0, 0; outIdx < outLen; inIdx, outIdx = inIdx+4, outIdx+2 {
		buf[outIdx] = (buf[inIdx] >> 2) + (buf[inIdx+2] >> 1) + (buf[inIdx+4] >> 2)
		buf[outIdx+1] = (buf[inIdx+1] >> 2) + (buf[inIdx+3] >> 1) + (buf[inIdx+5] >> 2)
	}
}

func (m *Micromod) tick() bool {
	if m.tickV-1 <= 0 {
		m.tickV = m.speed
		m.row()
	} else {
		for _, ch := range m.channels {
			ch.Tick()
		}
	}
	return m.playCount[m.seqPos][m.rowV] > 1
}

func (m *Micromod) row() {
	if m.nextRow < 0 {
		m.breakSeqPos = m.seqPos + 1
		m.nextRow = 0
	}
	if m.breakSeqPos >= 0 {
		if m.breakSeqPos >= m.module.SequenceLength() {
			m.breakSeqPos = 0
			m.nextRow = 0
		}
		m.seqPos = m.breakSeqPos
		for _, ch := range m.channels {
			ch.PlRow = 0
		}
		m.breakSeqPos = -1
	}
	m.rowV = m.nextRow
	count := m.playCount[m.seqPos][m.rowV]
	if m.plCount < 0 && count < 127 {
		m.playCount[m.seqPos][m.rowV] = count + 1
	}
	m.nextRow = m.rowV + 1
	if m.nextRow >= pattern.NumRows {
		m.nextRow = -1
	}
	for chIdx := 0; chIdx < len(m.channels); chIdx++ {
		ch := m.channels[chIdx]
		m.module.Pattern(m.module.SequenceEntry(m.seqPos)).Note(m.rowV, chIdx, m.note) // todo NB: some side effect maybe
		effect := m.note.Effect & 0xFF
		param := m.note.Parameter & 0xFF
		if effect == 0xE {
			effect = 0x10 | (param >> 4)
			param &= 0xF
		}
		if effect == 0 && param > 0 {
			effect = 0xE
		}
		ch.Row(m.note.Key, m.note.Instrument, effect, param)
		switch effect {
		case 0xB: // Pattern Jump
			if m.plCount < 0 {
				m.breakSeqPos = param
				m.nextRow = 0
			}
		case 0xD: // Pattern Break
			if m.plCount < 0 {
				if m.breakSeqPos < 0 {
					m.breakSeqPos = m.seqPos + 1
				}
				m.nextRow = (param>>4)&10 + (param & 0xF)
				if m.nextRow >= 64 {
					m.nextRow = 0
				}
			}
		case 0xF: // Set Speed
			if param > 0 {
				m.tickV = param
				m.speed = param
			} else {
				m.tempo = param
			}
		case 0x16: // Pattern Loop
			if param == 0 { // Set loop marker on this channel
				ch.PlRow = m.rowV
			}
			if ch.PlRow < m.rowV && m.breakSeqPos < 0 { // Marker valid
				if m.plCount < 0 { // Not already looping, begin
					m.plCount = param
					m.plChannel = chIdx
				}
				if m.plChannel == chIdx { // Next Loop
					if m.plCount == 0 { // Loop finished
						ch.PlRow = m.rowV + 1
					} else { // Loop
						m.nextRow = ch.PlRow
					}
					m.plCount -= 1
				}
			}
		case 0x1E: // Pattern Delay
			m.tickV = m.speed + m.speed*param
		}
	}
}

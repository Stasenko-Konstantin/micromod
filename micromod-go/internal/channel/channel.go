package channel

import (
	"micromod/internal/instrument"
	"micromod/internal/module"
)

var sineTable = []uint8{
	0, 24, 49, 74, 97, 120, 141, 161, 180, 197, 212, 224, 235, 244, 250, 253,
	255, 253, 250, 244, 235, 224, 212, 197, 180, 161, 141, 120, 97, 74, 49, 24,
}

type Channel struct {
	module                                                *module.Module
	noteKey, noteEffect, noteParam                        int
	noteIns, instrument, assigned                         int
	sampleOffset, sampleIdx, sampleFra                    int
	volume, panning, fineTune, freq, ampl                 int
	period, portaPeriod, portaSpeed, fxCount              int
	vibratoType, vibratoPhase, vibratoSpeed, vibratoDepth int
	tremoloType, tremoloPhase, tremoloSpeed, tremoloDepth int
	tremoloAdd, vibratoAdd, arpeggioAdd                   int
	id, randomSeed                                        int
	PlRow                                                 int
}

func New(module *module.Module, id int) *Channel {
	ch := &Channel{module: module, id: id}
	switch id & 0x3 {
	case 0, 3:
		ch.panning = 51
	case 1, 2:
		ch.panning = 204
	}
	ch.randomSeed = (id + 1) * 0xABCDEF
	return ch
}

func (c *Channel) Resample(mixBuf []int, offset, count, sampleRate int, interpolation bool) {
	if c.instrument > 0 && c.ampl > 0 {
		leftGain := (c.ampl * c.panning) >> 8
		rightGain := (c.ampl * (255 - c.panning)) >> 8
		step := (c.freq << (instrument.FP_SHIFT - 3)) / (sampleRate >> 3)
		c.module.Instrument(c.instrument).
			Audio(c.sampleIdx, c.sampleFra, step, leftGain, rightGain, mixBuf, offset, count, interpolation)
	}
}

func (c *Channel) UpdateSampleIdx(length, sampleRate int) {
	if c.instrument > 0 {
		step := (c.freq << (instrument.FP_SHIFT - 3)) / (sampleRate >> 3)
		c.sampleFra += step * length
		c.sampleIdx += c.sampleFra >> instrument.FP_SHIFT
		c.sampleIdx = c.module.Instrument(c.instrument).NormalizeSampleIdx(c.sampleIdx)
	}
}

func (c *Channel) Row(key, ins, effect, param int) {
	c.noteKey = key
	c.noteIns = ins
	c.noteEffect = effect
	c.noteParam = param
	if !(effect == 0x1D && param > 0) {
		// Not note delay
		c.trigger()
	}
	switch effect {
	case 0x3: // Tone Portamento
		if param > 0 {
			c.portaSpeed = param
		}
	case 0x4: // Vibrato
		if (param & 0xF0) > 0 {
			c.vibratoSpeed = param >> 4
		}
		if (param & 0x0F) > 0 {
			c.vibratoDepth = param & 0xF
		}
		c.vibrato()
	case 0x6: // Vebrato + Volume Slide
		c.vibrato()
	case 0x7: // Tremolo
		if (param & 0xF0) > 0 {
			c.tremoloSpeed = param >> 4
		}
		if (param & 0x0F) > 0 {
			c.tremoloDepth = param & 0xF
		}
		c.tremolo()
	case 0x8: // Set Panning. Not for 4-channel ProTracket
		if c.module.NumChannels() != 4 {
			if param < 128 {
				c.panning = param << 1
			} else {
				c.panning = 255
			}
		}
	case 0xC: // Set Volume
		if param > 64 {
			c.volume = 64
		} else {
			c.volume = param
		}
	case 0x11: // Fine Portamento Up
		c.period -= param
		if c.period < 0 {
			c.period = 0
		}
	case 0x12: // Fine Portamento Down
		c.period += param
		if c.period > 65535 {
			c.period = 65535
		}
	case 0x14: // Vibrato Waveform
		if param < 8 {
			c.vibratoType = param
		}
	case 0x17: // Tremolo Waveform
		if param < 8 {
			c.tremoloType = param
		}
	case 0x1A: // Fine Volume Up
		c.volume += param
		if c.volume > 64 {
			c.volume = 64
		}
	case 0x1B: // Fine Volume Down
		c.volume -= param
		if c.volume < 0 {
			c.volume = 0
		}
	case 0x1C: // Note Cut
		if param <= 0 {
			c.volume = 0
		}
	}
	c.updateFrequency()
}

func (c *Channel) Tick() {
	// todo
}

func (c *Channel) updateFrequency() {
	// todo
}

func (c *Channel) trigger() {
	// todo
}

func (c *Channel) vibrato() {
	// todo
}

func (c *Channel) tremolo() {
	// todo
}

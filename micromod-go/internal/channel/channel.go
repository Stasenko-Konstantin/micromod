package channel

import (
	"micromod/internal/module"
)

type Channel struct {
	module *module.Module
	// todo
	id, randomSeed int
	PlRow          int
}

func New(module *module.Module, id int) *Channel {
	// todo, randomSeed
	return &Channel{module: module, id: id}
}

func (c *Channel) Resample(mixBud []int, offset, count, sampleRate int, interpolation bool) {
	// todo
}

func (c *Channel) UpdateSampleIdx(length, sampleRate int) {
	// todo
}

func (c *Channel) Tick() {
	// todo
}

func (c *Channel) Row(key, ins, effect, param int) {
	// todo
}

package micromod

import (
	"micromod/internal/module"
	"micromod/internal/player"
)

type Micromod struct {
	module *module.Module
	player *player.Player
}

func New(modFilepath string, interpolation, loop bool) *Micromod {
	return &Micromod{module.New(modFilepath), player.New(interpolation, loop)}
}

func (m *Micromod) Run(clCh chan struct{}) {
	clCh <- struct{}{}
}

func (m *Micromod) GetModuleInfo() string {
	return m.module.GetModuleInfo()
}

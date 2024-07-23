package player

type Player struct {
	interpolation bool
	loop          bool
}

func New(interpolation, loop bool) *Player {
	return &Player{interpolation, loop}
}

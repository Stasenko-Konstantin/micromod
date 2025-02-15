package note

type Note struct {
	Key, Instrument, Effect, Parameter int
}

func New() *Note {
	return &Note{}
}

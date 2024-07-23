package module

type Module struct {
	modFilepath string
}

func New(modFilepath string) *Module {
	return &Module{modFilepath}
}

// todo
func (m *Module) GetModuleInfo() string {
	return ""
}

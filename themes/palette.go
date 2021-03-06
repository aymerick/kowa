package themes

// Palette represents a theme palette
type Palette struct {
	Name string            `json:"name"`
	Vars map[string]string `json:"vars"`
}

// NewPalette instanciates a new Palette
func NewPalette(name string) *Palette {
	return &Palette{
		Name: name,
		Vars: map[string]string{},
	}
}

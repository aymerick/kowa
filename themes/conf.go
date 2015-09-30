package themes

// Conf represents a theme configuration
type Conf struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Author   string     `json:"author,omitempty"`
	Homepage string     `json:"homepage,omitempty"`
	Palettes []*Palette `json:"palettes,omitempty"`
}

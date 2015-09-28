package themes

// Conf represents a theme configuration
type Conf struct {
	Name     string
	Author   Author
	Palettes []*Palette
}

// Author represents a theme author
type Author struct {
	Name     string
	Homepage string
}

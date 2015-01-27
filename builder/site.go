package builder

import "github.com/aymerick/kowa/models"

// Site vars for templates
type SiteVars struct {
	Name string // Site name
	Logo string // Site logo

	NavBar []*SiteNavBarItem // Navigation bar

	builder *SiteBuilder
}

type SiteNavBarItem struct {
	Url   string // Item URL
	Title string // Item title
}

func NewSiteVars(siteBuilder *SiteBuilder) *SiteVars {
	return &SiteVars{
		builder: siteBuilder,
	}
}

func NewSiteNavBarItem(url string, title string) *SiteNavBarItem {
	return &SiteNavBarItem{
		Url:   url,
		Title: title,
	}
}

func (vars *SiteVars) Fill() {
	site := vars.builder.site

	vars.Name = site.Name

	if logo := site.FindLogo(); logo != nil {
		vars.Logo = vars.builder.AddImage(logo, models.MEDIUM_KIND)
	}

	vars.NavBar = computeNavBarItems(vars.builder)
}

func computeNavBarItems(builder *SiteBuilder) []*SiteNavBarItem {
	result := make([]*SiteNavBarItem, 0)

	nodes := builder.NavBarNodes()
	for _, node := range nodes {
		result = append(result, NewSiteNavBarItem(node.Url(), node.Title))
	}

	return result
}

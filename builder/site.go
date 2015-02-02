package builder

import (
	"sort"

	"github.com/aymerick/kowa/models"
)

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
	Order int    // Item order
}

type NavBarItemsByOrder []*SiteNavBarItem

func NewSiteVars(siteBuilder *SiteBuilder) *SiteVars {
	return &SiteVars{
		builder: siteBuilder,
	}
}

// Fill site variables
func (vars *SiteVars) fill() {
	site := vars.builder.site

	vars.Name = site.Name

	if logo := site.FindLogo(); logo != nil {
		vars.Logo = vars.builder.addImage(logo, models.MEDIUM_KIND)
	}

	vars.NavBar = computeNavBarItems(vars.builder)
}

func computeNavBarItems(builder *SiteBuilder) []*SiteNavBarItem {
	result := []*SiteNavBarItem{}

	nodes := builder.navBarNodes()
	for _, node := range nodes {
		result = append(result, NewSiteNavBarItem(node.Url, node.Title, node.NavBarOrder))
	}

	// sort
	sort.Sort(NavBarItemsByOrder(result))

	return result
}

func NewSiteNavBarItem(url string, title string, order int) *SiteNavBarItem {
	return &SiteNavBarItem{
		Url:   url,
		Title: title,
		Order: order,
	}
}

// Implements sort.Interface
func (items NavBarItemsByOrder) Len() int {
	return len(items)
}

// Implements sort.Interface
func (items NavBarItemsByOrder) Swap(i, j int) {
	items[i], items[j] = items[j], items[i]
}

// Implements sort.Interface
func (items NavBarItemsByOrder) Less(i, j int) bool {
	return items[i].Order < items[j].Order
}

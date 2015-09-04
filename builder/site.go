package builder

import (
	"sort"

	"github.com/nicksnyder/go-i18n/i18n"
)

// SiteVars holds site theme variables
type SiteVars struct {
	Name         string
	Logo         *ImageVars
	Tagline      string
	NameInNavBar bool
	Favicon      bool

	Membership string

	BaseUrl  string
	BasePath string

	Facebook   string
	Twitter    string
	GooglePlus string

	GoogleAnalytics string

	NavBar []*SiteNavBarItem // Navigation bar

	builder *SiteBuilder
}

// SiteNavBarItem represents an item in navigation bar
type SiteNavBarItem struct {
	Url   string // Item URL
	Title string // Item title
	Order int    // Item order
}

// NavBarItemsByOrder holds an ordered list of navigation bar items
type NavBarItemsByOrder []*SiteNavBarItem

// NewSiteVars instanciates a new SiteVars
func NewSiteVars(siteBuilder *SiteBuilder) *SiteVars {
	return &SiteVars{
		builder: siteBuilder,
	}
}

// Fill site variables
func (vars *SiteVars) fill() {
	site := vars.builder.site

	T := i18n.MustTfunc(site.Lang)

	name := site.Name
	if name == "" {
		name = T("empty_site_name")
	}

	vars.Name = name
	vars.BaseUrl = site.BaseUrl()
	vars.BasePath = vars.builder.basePath()
	vars.Tagline = site.Tagline
	vars.NameInNavBar = site.NameInNavBar

	if img := site.FindFavicon(); img != nil {
		vars.Favicon = true
	}

	if file := site.FindMembership(); file != nil {
		vars.Membership = vars.builder.addFile(file)
	}

	vars.Facebook = site.Facebook
	vars.Twitter = site.Twitter
	vars.GooglePlus = site.GooglePlus

	vars.GoogleAnalytics = site.GoogleAnalytics

	if logo := site.FindLogo(); logo != nil {
		vars.Logo = vars.builder.addImage(logo)
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

// NewSiteNavBarItem instanciates a new SiteNavBarItem
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

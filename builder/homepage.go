package builder

type HomepageBuilder struct {
}

func NewHomepageBuilder() *HomepageBuilder {
	return &HomepageBuilder{}
}

func (builder *HomepageBuilder) Fill(page *SitePage, site *Site) error {
	page.Title = "Homepage"

	page.Meta = &SitePageMeta{
		Description: "Homepage test page",
	}

	page.Content = "Soon"

	return nil
}

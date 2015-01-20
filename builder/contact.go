package builder

type ContactBuilder struct {
}

func NewContactBuilder() *ContactBuilder {
	return &ContactBuilder{}
}

func (builder *ContactBuilder) Fill(page *SitePage, site *Site) error {
	page.Title = "Contact"

	page.Meta = &SitePageMeta{
		Description: "Contact test page",
	}

	page.Content = "Soon"

	return nil
}

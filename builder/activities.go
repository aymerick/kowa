package builder

type ActivitiesBuilder struct {
}

func NewActivitiesBuilder() *ActivitiesBuilder {
	return &ActivitiesBuilder{}
}

func (builder *ActivitiesBuilder) Fill(page *SitePage, site *Site) error {
	page.Title = "Activities"
	page.Meta = &SitePageMeta{
		Description: "Activities test page",
	}
	page.BodyClass = "activities"
	page.Content = []string{"one", "two", "three<br />", "four"}

	return nil
}

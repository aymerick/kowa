package builder

import (
	"html/template"
	"time"

	"github.com/aymerick/kowa/models"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// Builder for activities page
type ActivitiesBuilder struct {
	*NodeBuilderBase

	// loaded activities
	activitiesContents []*ActivityContent
}

// Activity content for template
type ActivityContent struct {
	Date    time.Time
	Cover   string
	Title   string
	Summary template.HTML
	Body    template.HTML
}

// Activities content for template
type ActivitiesContent struct {
	Title   string
	Tagline string

	Activities []*ActivityContent
}

func init() {
	RegisterNodeBuilder(KIND_ACTIVITIES, NewActivitiesBuilder)
}

func NewActivitiesBuilder(siteBuilder *SiteBuilder) NodeBuilder {
	return &ActivitiesBuilder{
		NodeBuilderBase: &NodeBuilderBase{
			nodeKind:    KIND_ACTIVITIES,
			siteBuilder: siteBuilder,
		},
	}
}

// NodeBuilder
func (builder *ActivitiesBuilder) Load() {
	// fetch activities
	activitiesContents := builder.activities()
	if len(activitiesContents) > 0 {
		// build activities page
		node := builder.newNode()
		node.fillUrl(node.Kind)

		title := "Activities"
		tagline := "" // @todo

		node.Title = title
		node.Meta = &NodeMeta{Description: tagline}
		node.Content = &ActivitiesContent{
			Title:      title,
			Tagline:    tagline,
			Activities: activitiesContents,
		}
		node.InNavBar = true
		node.NavBarOrder = 1

		builder.addNode(node)
	}
}

// NodeBuilder
func (builder *ActivitiesBuilder) Data(name string) interface{} {
	switch name {
	case "activities":
		return builder.activities()
	}

	return nil
}

// returns activities contents
func (builder *ActivitiesBuilder) activities() []*ActivityContent {
	if len(builder.activitiesContents) == 0 {
		// fetch activities
		for _, activity := range *builder.site().FindAllActivities() {
			activityContent := builder.NewActivityContent(activity)

			builder.activitiesContents = append(builder.activitiesContents, activityContent)
		}
	}

	return builder.activitiesContents
}

func (builder *ActivitiesBuilder) NewActivityContent(activity *models.Activity) *ActivityContent {
	result := &ActivityContent{
		Date:  activity.CreatedAt,
		Title: activity.Title,
	}

	cover := activity.FindCover()
	if cover != nil {
		result.Cover = builder.addImage(cover, models.MEDIUM_KIND)
	}

	htmlSummary := blackfriday.MarkdownCommon([]byte(activity.Summary))
	result.Summary = template.HTML(bluemonday.UGCPolicy().SanitizeBytes(htmlSummary))

	htmlBody := blackfriday.MarkdownCommon([]byte(activity.Body))
	result.Body = template.HTML(bluemonday.UGCPolicy().SanitizeBytes(htmlBody))

	return result
}

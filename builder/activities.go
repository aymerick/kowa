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
	Date  time.Time     // CreatedAt
	Cover string        // Cover URL
	Title string        // Title
	Body  template.HTML // Body
}

// Activities content for template
type ActivitiesContent struct {
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

	// build activities page
	node := builder.newNode()
	node.fillUrl(node.Kind)

	node.Title = "Activities"
	node.Meta = &NodeMeta{Description: "Activities test page"}
	node.Content = NewActivitiesContent(activitiesContents)
	node.InNavBar = true

	builder.addNode(node)
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
		for _, activity := range *builder.SiteBuilder().site.FindAllActivities() {
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

	html := blackfriday.MarkdownCommon([]byte(activity.Body))
	result.Body = template.HTML(bluemonday.UGCPolicy().SanitizeBytes(html))

	return result
}

func NewActivitiesContent(activities []*ActivityContent) *ActivitiesContent {
	return &ActivitiesContent{
		Activities: activities,
	}
}

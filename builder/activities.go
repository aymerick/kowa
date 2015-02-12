package builder

import (
	"html/template"
	"time"

	"github.com/aymerick/kowa/models"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// Activities node builder
type ActivitiesBuilder struct {
	*NodeBuilderBase

	// loaded activities
	activitiesVars []*ActivityVars
}

// Activities node content
type ActivitiesContent struct {
	Node *Node

	Activities []*ActivityVars
}

// Activity vars
type ActivityVars struct {
	Date    time.Time
	Cover   string
	Title   string
	Summary template.HTML
	Body    template.HTML
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
	activitiesVars := builder.activities()
	if len(activitiesVars) > 0 {
		// build activities page
		node := builder.newNode()
		node.fillUrl(node.Kind)

		title := "Activities" // @todo i18n
		tagline := ""         // @todo Fill

		node.Title = title
		node.Tagline = tagline
		node.Meta = &NodeMeta{Description: tagline}
		node.InNavBar = true
		node.NavBarOrder = 1

		node.Content = &ActivitiesContent{
			Node:       node,
			Activities: activitiesVars,
		}

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
func (builder *ActivitiesBuilder) activities() []*ActivityVars {
	if len(builder.activitiesVars) == 0 {
		// fetch activities
		for _, activity := range *builder.site().FindAllActivities() {
			activityVars := builder.NewActivityVars(activity)

			builder.activitiesVars = append(builder.activitiesVars, activityVars)
		}
	}

	return builder.activitiesVars
}

func (builder *ActivitiesBuilder) NewActivityVars(activity *models.Activity) *ActivityVars {
	result := &ActivityVars{
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

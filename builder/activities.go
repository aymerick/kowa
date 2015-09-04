package builder

import (
	"time"

	"github.com/nicksnyder/go-i18n/i18n"

	"github.com/aymerick/kowa/models"
	"github.com/aymerick/raymond"
)

// ActivitiesBuilder builds activities page
type ActivitiesBuilder struct {
	*NodeBuilderBase

	// loaded activities
	activitiesVars []*ActivityVars
}

// ActivitiesContent represents activities node content
type ActivitiesContent struct {
	Activities []*ActivityVars
}

// ActivityVars represents activity vars
type ActivityVars struct {
	Date    time.Time
	Cover   *ImageVars
	Title   string
	Summary raymond.SafeString
	Body    raymond.SafeString
}

func init() {
	RegisterNodeBuilder(kindActivities, NewActivitiesBuilder)
}

// NewActivitiesBuilder instanciates a new NodeBuilder
func NewActivitiesBuilder(siteBuilder *SiteBuilder) NodeBuilder {
	return &ActivitiesBuilder{
		NodeBuilderBase: &NodeBuilderBase{
			nodeKind:    kindActivities,
			siteBuilder: siteBuilder,
		},
	}
}

// Load is part of NodeBuilder interface
func (builder *ActivitiesBuilder) Load() {
	// fetch activities
	activitiesVars := builder.activities()
	if len(activitiesVars) == 0 {
		return
	}

	// get page settings
	title, tagline, cover, disabled := builder.pageSettings(models.PageKindActivities)
	if disabled {
		return
	}

	T := i18n.MustTfunc(builder.siteLang())
	slug := T("activities")

	if title == "" {
		title = slug
	}

	// build node
	node := builder.newNode()
	node.fillURL(slug)

	node.Title = title
	node.Tagline = tagline
	node.Cover = cover

	node.Meta = &NodeMeta{Description: tagline}

	node.InNavBar = true
	node.NavBarOrder = 1

	node.Content = &ActivitiesContent{
		Activities: activitiesVars,
	}

	builder.addNode(node)
}

// Data is part of NodeBuilder interface
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

// NewActivityVars instanciates a new ActivityVars
func (builder *ActivitiesBuilder) NewActivityVars(activity *models.Activity) *ActivityVars {
	result := &ActivityVars{
		Date:  activity.CreatedAt,
		Title: activity.Title,
	}

	cover := activity.FindCover()
	if cover != nil {
		result.Cover = builder.addImage(cover)
	}

	result.Summary = generateHTML(models.FormatHTML, activity.Summary)
	result.Body = generateHTML(models.FormatHTML, activity.Body)

	return result
}

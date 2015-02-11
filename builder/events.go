package builder

import (
	"html/template"
	"path"

	"github.com/aymerick/kowa/models"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// Builder for events pages
type EventsBuilder struct {
	*NodeBuilderBase

	events []*EventNodeContentPair
}

// Event content for template
type EventContent struct {
	StartDate string
	EndDate   string
	Cover     string
	Title     string
	Place     string
	Body      template.HTML
	Url       string
}

// Event with associated Node Content
type EventNodeContentPair struct {
	event       *models.Event
	nodeContent *EventContent
}

// Event list content for template
type EventListContent struct {
	Title   string
	Tagline string

	Events []*EventContent
	// PrevPage string
	// NextPage string
}

func init() {
	RegisterNodeBuilder(KIND_EVENTS, NewEventsBuilder)
}

// Instanciate a new builder
func NewEventsBuilder(siteBuilder *SiteBuilder) NodeBuilder {
	return &EventsBuilder{
		NodeBuilderBase: &NodeBuilderBase{
			nodeKind:    KIND_EVENT,
			siteBuilder: siteBuilder,
		},
	}
}

func NewEventNodeContentPair(event *models.Event, nodeContent *EventContent) *EventNodeContentPair {
	return &EventNodeContentPair{
		event:       event,
		nodeContent: nodeContent,
	}
}

// NodeBuilder
func (builder *EventsBuilder) Load() {
	builder.loadEvents()
	builder.loadEventsLists()
}

// Build all events
func (builder *EventsBuilder) loadEvents() {
	for _, event := range *builder.site().FindAllEvents() {
		builder.loadEvent(event)
	}
}

// Build event page
func (builder *EventsBuilder) loadEvent(event *models.Event) {
	node := builder.newNode()
	node.fillUrl(path.Join("events", event.Slug())) // @todo i18n

	node.Title = event.Title
	node.Meta = &NodeMeta{
		Description: "", // @todo !!!
	}

	eventContent := builder.NewEventContent(event, node)

	node.Content = eventContent

	builder.addNode(node)

	builder.events = append(builder.events, NewEventNodeContentPair(event, eventContent))
}

// Instanciate a new event content
func (builder *EventsBuilder) NewEventContent(event *models.Event, node *Node) *EventContent {
	result := &EventContent{
		StartDate: event.StartDate.Format("02/01/06 15:04"),
		EndDate:   event.EndDate.Format("02/01/06 15:04"),
		Title:     event.Title,
		Place:     event.Place,
		Url:       node.Url,
	}

	cover := event.FindCover()
	if cover != nil {
		result.Cover = builder.addImage(cover, models.MEDIUM_KIND)
	}

	html := blackfriday.MarkdownCommon([]byte(event.Body))
	result.Body = template.HTML(bluemonday.UGCPolicy().SanitizeBytes(html))

	return result
}

// Build events list pages
func (builder *EventsBuilder) loadEventsLists() {
	if len(builder.events) > 0 {
		// @todo pagination
		node := builder.newNodeForKind(KIND_EVENTS)
		node.fillUrl(KIND_EVENTS)

		title := "Events"
		tagline := "" // @todo

		node.Title = title
		node.Meta = &NodeMeta{Description: tagline}
		node.Content = &EventListContent{
			Title:   title,
			Tagline: tagline,
			Events:  computesEventContents(builder.events),
		}
		node.InNavBar = true
		node.NavBarOrder = 10

		builder.addNode(node)
	}
}

func computesEventContents(events []*EventNodeContentPair) []*EventContent {
	eventContents := []*EventContent{}

	for _, eventNodeContent := range events {
		eventContents = append(eventContents, eventNodeContent.nodeContent)
	}

	return eventContents
}

package builder

import (
	"fmt"
	"html/template"
	"path"
	"sort"
	"time"

	"github.com/aymerick/kowa/models"
	"github.com/aymerick/kowa/utils"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// Event nodes builder
type EventsBuilder struct {
	*NodeBuilderBase

	events     []*EventContent
	pastEvents []*EventContent
}

// Event node content
type EventContent struct {
	Node  *Node
	Model *models.Event

	Cover *ImageVars
	Title string
	Place string
	Body  template.HTML
	Url   string

	Dates string

	StartDateRFC3339 string
	StartDateTime    string
	StartDate        string
	StartWeekday     string
	StartDay         string
	StartMonth       string
	StartYear        string
	StartTime        string

	EndDateRFC3339 string
	EndDateTime    string
	EndDate        string
	EndWeekday     string
	EndDay         string
	EndMonth       string
	EndYear        string
	EndTime        string
}

// Sortable event node contents
type EventContentsByStartDate []*EventContent

// Events node content
type EventsContent struct {
	Node *Node

	Events     []*EventContent
	PastEvents []*EventContent

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

// Computes event slug
func eventSlug(event *models.Event) string {
	year, month, day := event.StartDate.Date()

	title := event.Title
	if len(title) > MAX_SLUG {
		title = title[:MAX_SLUG]
	}

	return fmt.Sprintf("%d/%02d/%02d/%s", year, month, day, utils.Urlify(title))
}

// Build event page
func (builder *EventsBuilder) loadEvent(event *models.Event) {
	node := builder.newNode()
	node.fillUrl(path.Join("events", eventSlug(event))) // @todo i18n

	title := "Events" // @todo i18n
	tagline := ""     // @todo Fill

	node.Title = title
	node.Tagline = tagline
	node.Meta = &NodeMeta{
		Title:       fmt.Sprintf("%s - %s", event.Title, builder.site().Name),
		Description: tagline,
	}

	eventContent := builder.NewEventContent(event, node)

	node.Content = eventContent

	builder.addNode(node)

	if time.Now().After(event.EndDate) {
		builder.pastEvents = append(builder.pastEvents, eventContent)
	} else {
		builder.events = append(builder.events, eventContent)
	}
}

// Instanciate a new event content
func (builder *EventsBuilder) NewEventContent(event *models.Event, node *Node) *EventContent {
	result := &EventContent{
		Node:  node,
		Model: event,

		Title: event.Title,
		Place: event.Place,
		Url:   node.Url,

		StartDateRFC3339: event.StartDate.Format(time.RFC3339),
		// @todo i18n
		StartDateTime: event.StartDate.Format("Mon Jan 02 3:04PM"),
		StartDate:     event.StartDate.Format("Mon Jan 02"),
		StartWeekday:  event.StartDate.Format("Mon"),
		StartDay:      event.StartDate.Format("02"),
		StartMonth:    event.StartDate.Format("Jan"),
		StartYear:     event.StartDate.Format("2006"),
		StartTime:     event.StartDate.Format("3:04PM"),

		EndDateRFC3339: event.EndDate.Format(time.RFC3339),
		// @todo i18n
		EndDateTime: event.EndDate.Format("Mon Jan 02 3:04PM"),
		EndDate:     event.EndDate.Format("Mon Jan 02"),
		EndWeekday:  event.EndDate.Format("Mon"),
		EndDay:      event.EndDate.Format("02"),
		EndMonth:    event.EndDate.Format("Jan"),
		EndYear:     event.EndDate.Format("2006"),
		EndTime:     event.EndDate.Format("3:04PM"),
	}

	if result.StartDate == result.EndDate {
		// @todo i18n
		result.Dates = fmt.Sprintf("%s from %s to %s", result.StartDate, result.StartTime, result.EndTime)
	} else {
		// @todo i18n
		result.Dates = fmt.Sprintf("From %s to %s", result.StartDateTime, result.EndDateTime)
	}

	cover := event.FindCover()
	if cover != nil {
		result.Cover = builder.addImage(cover)
	}

	html := blackfriday.MarkdownCommon([]byte(event.Body))
	result.Body = template.HTML(bluemonday.UGCPolicy().SanitizeBytes(html))

	return result
}

// Build events list pages
func (builder *EventsBuilder) loadEventsLists() {
	if len(builder.events) > 0 || len(builder.pastEvents) > 0 {
		// @todo pagination
		node := builder.newNodeForKind(KIND_EVENTS)
		node.fillUrl(KIND_EVENTS)

		title := "Events" // @todo i18n
		tagline := ""     // @todo Fill

		node.Title = title
		node.Tagline = tagline
		node.Meta = &NodeMeta{Description: tagline}
		node.InNavBar = true
		node.NavBarOrder = 10

		events := builder.events
		sort.Sort(EventContentsByStartDate(events))

		pastEvents := builder.pastEvents
		sort.Sort(sort.Reverse(EventContentsByStartDate(pastEvents)))
		if len(pastEvents) > MAX_PAST_EVENTS {
			pastEvents = pastEvents[:MAX_PAST_EVENTS]
		}

		node.Content = &EventsContent{
			Node:       node,
			Events:     events,
			PastEvents: pastEvents,
		}

		builder.addNode(node)
	}
}

//
// EventContentsByStartDate
//

// Implements sort.Interface
func (events EventContentsByStartDate) Len() int {
	return len(events)
}

// Implements sort.Interface
func (events EventContentsByStartDate) Swap(i, j int) {
	events[i], events[j] = events[j], events[i]
}

// Implements sort.Interface
func (events EventContentsByStartDate) Less(i, j int) bool {
	return events[i].Model.StartDate.Before(events[j].Model.StartDate)
}

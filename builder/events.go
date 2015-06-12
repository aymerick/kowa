package builder

import (
	"fmt"
	"path"
	"sort"
	"time"

	"github.com/nicksnyder/go-i18n/i18n"

	"github.com/aymerick/kowa/models"
	"github.com/aymerick/raymond"
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
	Body  raymond.SafeString
	Url   string

	Dates string

	StartDateRFC3339  string
	StartDateTime     string
	StartDate         string
	StartWeekday      string
	StartWeekdayShort string
	StartDay          string
	StartMonth        string
	StartMonthShort   string
	StartYear         string
	StartTime         string

	EndDateRFC3339  string
	EndDateTime     string
	EndDate         string
	EndWeekday      string
	EndWeekdayShort string
	EndDay          string
	EndMonth        string
	EndMonthShort   string
	EndYear         string
	EndTime         string
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

	return fmt.Sprintf("%d/%02d/%02d/%s", year, month, day, title)
}

// Build event page
func (builder *EventsBuilder) loadEvent(event *models.Event) {
	// get page settings
	title, tagline, cover, disabled := builder.pageSettings(models.PAGE_KIND_EVENTS)
	if disabled {
		return
	}

	T := i18n.MustTfunc(builder.siteLang())
	slug := T("events")

	if title == "" {
		title = slug
	}

	// build node
	node := builder.newNode()
	node.fillUrl(path.Join(slug, eventSlug(event)))

	node.Title = title
	node.Tagline = tagline

	node.Meta = &NodeMeta{
		Title:       fmt.Sprintf("%s - %s", event.Title, builder.site().Name),
		Description: tagline,
		Type:        "article",
	}

	eventContent := builder.NewEventContent(event, node)
	node.Content = eventContent

	if eventContent.Cover != nil {
		node.Cover = eventContent.Cover
	} else {
		node.Cover = cover
	}

	builder.addNode(node)

	if time.Now().After(event.EndDate) {
		builder.pastEvents = append(builder.pastEvents, eventContent)
	} else {
		builder.events = append(builder.events, eventContent)
	}
}

// Instanciate a new event content
func (builder *EventsBuilder) NewEventContent(event *models.Event, node *Node) *EventContent {
	T := i18n.MustTfunc(builder.siteLang())

	result := &EventContent{
		Node:  node,
		Model: event,

		Title: event.Title,
		Place: event.Place,
		Url:   node.Url,

		StartDateRFC3339:  event.StartDate.Format(time.RFC3339),
		StartWeekday:      T("weekday_" + event.StartDate.Format("Monday")),
		StartWeekdayShort: T("weekday_short_" + event.StartDate.Format("Mon")),
		StartDay:          event.StartDate.Format("02"),
		StartMonth:        T("month_" + event.StartDate.Format("January")),
		StartMonthShort:   T("month_short_" + event.StartDate.Format("Jan")),
		StartYear:         event.StartDate.Format("2006"),
		StartTime:         event.StartDate.Format(T("format_time")),

		EndDateRFC3339:  event.EndDate.Format(time.RFC3339),
		EndWeekday:      T("weekday_" + event.EndDate.Format("Monday")),
		EndWeekdayShort: T("weekday_short_" + event.EndDate.Format("Mon")),
		EndDay:          event.EndDate.Format("02"),
		EndMonth:        T("month_" + event.EndDate.Format("January")),
		EndMonthShort:   T("month_short_" + event.EndDate.Format("Jan")),
		EndYear:         event.EndDate.Format("2006"),
		EndTime:         event.EndDate.Format(T("format_time")),
	}

	result.StartDateTime = T("event_format_datetime", map[string]interface{}{
		"Year":    result.StartYear,
		"Month":   result.StartMonth,
		"Day":     result.StartDay,
		"Time":    result.StartTime,
		"Weekday": result.StartWeekday,
	})

	result.StartDate = T("event_format_date", map[string]interface{}{
		"Year":    result.StartYear,
		"Month":   result.StartMonth,
		"Day":     result.StartDay,
		"Weekday": result.StartWeekday,
	})

	result.EndDateTime = T("event_format_datetime", map[string]interface{}{
		"Year":    result.EndYear,
		"Month":   result.EndMonth,
		"Day":     result.EndDay,
		"Time":    result.EndTime,
		"Weekday": result.EndWeekday,
	})

	result.EndDate = T("event_format_date", map[string]interface{}{
		"Year":    result.EndYear,
		"Month":   result.EndMonth,
		"Day":     result.EndDay,
		"Weekday": result.EndWeekday,
	})

	if result.StartDate == result.EndDate {
		result.Dates = T("date_times_interval", map[string]interface{}{
			"StartDate": result.StartDate,
			"StartTime": result.StartTime,
			"EndTime":   result.EndTime,
		})
	} else {
		result.Dates = T("dates_interval", map[string]interface{}{
			"StartDateTime": result.StartDateTime,
			"EndDateTime":   result.EndDateTime,
		})
	}

	cover := event.FindCover()
	if cover != nil {
		result.Cover = builder.addImage(cover)
	}

	result.Body = generateHTML(event.Format, event.Body)

	return result
}

// Build events list pages
// @todo pagination
func (builder *EventsBuilder) loadEventsLists() {
	if len(builder.events) == 0 && len(builder.pastEvents) == 0 {
		return
	}

	// get page settings
	title, tagline, cover, disabled := builder.pageSettings(models.PAGE_KIND_EVENTS)
	if disabled {
		return
	}

	T := i18n.MustTfunc(builder.siteLang())
	slug := T("events")

	if title == "" {
		title = slug
	}

	// build node
	node := builder.newNodeForKind(KIND_EVENTS)
	node.fillUrl(slug)

	node.Title = title
	node.Tagline = tagline
	node.Cover = cover

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

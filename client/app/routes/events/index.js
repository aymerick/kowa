import AuthenticatedRoute from 'kowa/routes/authenticated';

var EventsIndexRoute = AuthenticatedRoute.extend({
  // redirects /:site_id/events to /:site_id/events/:event_id
  beforeModel: function () {
    var events = this.modelFor('events');
    var firstEvent = events.get('firstObject');
    if (firstEvent) {
      this.transitionTo('events.event', firstEvent);
    }
  }
});

export default EventsIndexRoute;

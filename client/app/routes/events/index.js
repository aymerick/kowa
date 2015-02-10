import AuthenticatedRoute from 'kowa/routes/authenticated';

var EventIndexRoute = AuthenticatedRoute.extend({
  // redirects /:site_id/events to /:site_id/events/:event_id
  beforeModel: function () {
    var events = this.modelFor('events');
    var firstEvent = events.get('firstObject');
    if (firstEvent) {
      this.transitionTo('event', firstEvent);
    }
  }
});

export default EventIndexRoute;

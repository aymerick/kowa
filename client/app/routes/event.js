import AuthenticatedRoute from 'kowa/routes/authenticated';

var EventRoute = AuthenticatedRoute.extend({
  model: function(params) {
    return this.store.find('event', params.event_id);
  }
});

export default EventRoute;

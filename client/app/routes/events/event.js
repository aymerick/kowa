import AuthenticatedRoute from 'kowa/routes/authenticated';
import ContentEditionRoute from 'kowa/mixins/content-edition-route';

var EventsEventRoute = AuthenticatedRoute.extend(ContentEditionRoute, {
  model: function(params) {
    return this.store.find('event', params.event_id);
  }
});

export default EventsEventRoute;

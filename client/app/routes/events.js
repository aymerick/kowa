import AuthenticatedRoute from 'kowa/routes/authenticated';
import ContentListRoute from 'kowa/mixins/content-list-route';

var EventsRoute = AuthenticatedRoute.extend(ContentListRoute, {
  contentModelType: 'event'
});

export default EventsRoute;

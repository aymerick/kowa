import AuthenticatedRoute from 'kowa/routes/authenticated';
import ContentListRoute from 'kowa/mixins/content-list-route';
import StyleBodyMixin from 'kowa/mixins/style-body';

var EventsRoute = AuthenticatedRoute.extend(ContentListRoute, StyleBodyMixin, {
  contentModelType: 'event',
  classNames: ['events']
});

export default EventsRoute;

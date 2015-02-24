import AuthenticatedRoute from 'kowa/routes/authenticated';
import ContentListRouteMixin from 'kowa/mixins/content-list-route';
import StyleBodyMixin from 'kowa/mixins/style-body';

var EventsRoute = AuthenticatedRoute.extend(ContentListRouteMixin, StyleBodyMixin, {
  contentModelType: 'event',
  classNames: ['events']
});

export default EventsRoute;

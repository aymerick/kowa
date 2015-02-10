import AuthenticatedRoute from 'kowa/routes/authenticated';
import ContentListRoute from 'kowa/mixins/content-list-route';
import StyleBodyMixin from 'kowa/mixins/style-body';

var PagesRoute = AuthenticatedRoute.extend(ContentListRoute, StyleBodyMixin, {
  contentModelType: 'page',
  classNames: ['pages']
});

export default PagesRoute;

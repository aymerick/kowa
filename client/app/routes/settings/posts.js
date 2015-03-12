import AuthenticatedRoute from 'kowa/routes/authenticated';
import PageSettingsRoute from 'kowa/mixins/page-settings-route';

var SettingsPostsRoute = AuthenticatedRoute.extend(PageSettingsRoute, {
  pageKind: 'posts'
});

export default SettingsPostsRoute;

import AuthenticatedRoute from 'kowa/routes/authenticated';
import PageSettingsRoute from 'kowa/mixins/page-settings-route';

var SettingsActivitiesPageRoute = AuthenticatedRoute.extend(PageSettingsRoute, {
  pageKind: 'activities'
});

export default SettingsActivitiesPageRoute;

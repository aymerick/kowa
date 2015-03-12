import AuthenticatedRoute from 'kowa/routes/authenticated';
import PageSettingsRoute from 'kowa/mixins/page-settings-route';

var SettingsContactPageRoute = AuthenticatedRoute.extend(PageSettingsRoute, {
  pageKind: 'contact'
});

export default SettingsContactPageRoute;

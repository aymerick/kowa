import AuthenticatedRoute from 'kowa/routes/authenticated';
import PageSettingsRoute from 'kowa/mixins/page-settings-route';

var SettingsEventsRoute = AuthenticatedRoute.extend(PageSettingsRoute, {
  pageKind: 'events'
});

export default SettingsEventsRoute;

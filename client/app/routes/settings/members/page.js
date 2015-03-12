import AuthenticatedRoute from 'kowa/routes/authenticated';
import PageSettingsRoute from 'kowa/mixins/page-settings-route';

var SettingsMembersPageRoute = AuthenticatedRoute.extend(PageSettingsRoute, {
  pageKind: 'members'
});

export default SettingsMembersPageRoute;

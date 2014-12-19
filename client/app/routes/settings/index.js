import AuthenticatedRoute from 'kowa/routes/authenticated';

var SettingsIndexRoute = AuthenticatedRoute.extend({
  beforeModel: function () {
    // redirects /settings to to /settings/general
    this.transitionTo('settings.general');
  }
});

export default SettingsIndexRoute;

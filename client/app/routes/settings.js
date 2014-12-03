import AuthenticatedRoute from 'kowa/routes/authenticated';

var SettingsRoute = AuthenticatedRoute.extend({
  beforeModel: function () {
    this.transitionTo('settings.general');
  }
});

export default SettingsRoute;

import AuthenticatedRoute from 'kowa/routes/authenticated';

var SettingsGeneralRoute = AuthenticatedRoute.extend({
  model: function() {
    return this.modelFor('site');
  }
});

export default SettingsGeneralRoute;

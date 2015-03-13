import AuthenticatedRoute from 'kowa/routes/authenticated';

var SettingsAccountRoute = AuthenticatedRoute.extend({
  model: function() {
    return this.session.get('currentUser');
  }
});

export default SettingsAccountRoute;

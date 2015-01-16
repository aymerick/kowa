import AuthenticatedRoute from 'kowa/routes/authenticated';

var SettingsActivitiesActivityRoute = AuthenticatedRoute.extend({
  model: function(params) {
    return this.store.find('activity', params.activity_id);
  }
});

export default SettingsActivitiesActivityRoute;

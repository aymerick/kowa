import AuthenticatedRoute from 'kowa/routes/authenticated';

var SettingsActivitiesIndexRoute = AuthenticatedRoute.extend({
  // redirects /:site_id/settings/activities to /:site_id/settings/activities/:activity_id
  beforeModel: function () {
    var activities = this.modelFor('settings.activities');
    var firstActivity = activities.get('firstObject');
    if (firstActivity) {
      this.transitionTo('settings.activities.activity', firstActivity);
    }
  }
});

export default SettingsActivitiesIndexRoute;

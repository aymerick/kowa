import Ember from 'ember';
import AuthenticatedRoute from 'kowa/routes/authenticated';

var SettingsActivitiesIndexRoute = AuthenticatedRoute.extend({
  model: function() {
    var site = this.modelFor('site');

    // using filter() allows the template to auto-update when new models are pulled in from the server
    return this.store.filter('activity', { 'site': site.get('id') }, function () {
      // nothing to filter
      return true;
    });
  }
});

export default SettingsActivitiesIndexRoute;

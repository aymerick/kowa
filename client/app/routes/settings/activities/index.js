import AuthenticatedRoute from 'kowa/routes/authenticated';

var SettingsActivitiesIndexRoute = AuthenticatedRoute.extend({
  model: function() {
    var site = this.modelFor('site');

    // 1. Caching resultPromise permits to avoid sending a request to server everytime we transition into this route
    // 2. Using filter() allows the template to auto-update when new models are pulled in from the server
    var resultPromise = this.get('resultPromise') || this.store.filter('activity', { 'site': site.get('id') }, function (activity) {
      return activity.get('site.id') == site.get('id');
    });

    this.set('resultPromise', resultPromise);

    return resultPromise;
  }
});

export default SettingsActivitiesIndexRoute;

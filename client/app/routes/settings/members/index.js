import AuthenticatedRoute from 'kowa/routes/authenticated';

var SettingsMembersIndexRoute = AuthenticatedRoute.extend({
  model: function() {
    var site = this.modelFor('site');

    // 1. Caching resultPromise permits to avoid sending a request to server everytime we transition into this route
    // 2. Using filter() allows the template to auto-update when new models are pulled in from the server
    var resultPromise = this.get('resultPromise') || this.store.filter('member', { 'site': site.get('id') }, function () {
      // nothing to filter
      return true;
    });

    this.set('resultPromise', resultPromise);

    return resultPromise;
  }
});

export default SettingsMembersIndexRoute;

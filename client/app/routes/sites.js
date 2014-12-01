import AuthenticatedRoute from 'kowa/routes/authenticated';

var SitesRoute = AuthenticatedRoute.extend({
  model: function() {
    return this.session.get('currentUser').then(function(currentUser){
      return currentUser.get('sites');
    });
  },

  afterModel: function(sites, transition) {
    // redirects / to /:site_id
    if (!transition.params.site) {
      var firstSite = sites.get('firstObject');
      if (firstSite) {
        this.transitionTo('site', firstSite);
      }
    }
  }
});

export default SitesRoute;

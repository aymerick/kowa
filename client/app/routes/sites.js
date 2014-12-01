import AuthenticatedRoute from 'kowa/routes/authenticated';

var SitesRoute = AuthenticatedRoute.extend({
  model: function() {
    return this.session.get('currentUser').then(function(currentUser){
      return currentUser.get('sites');
    });
  },

  afterModel: function(sites, transition) {
    // catch when accessing "/" endpoint
    if (!transition.params.site) {
      // redirects to first site
      // @todo Handle case when user have no site yet !
      this.transitionTo('site', sites.get('firstObject'));
    }
  }
});

export default SitesRoute;

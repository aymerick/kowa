import AuthenticatedRoute from 'kowa/routes/authenticated';

var SitesRoute = AuthenticatedRoute.extend({
  model: function() {
    return this.session.get('currentUser').then(function(currentUser){
      return currentUser.get('sites');
    });
  },

  afterModel: function(sites, transition) {
    this.transitionTo('site', sites.get('firstObject'));
  }
});

export default SitesRoute;

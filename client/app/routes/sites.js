import AuthenticatedRoute from 'kowa/routes/authenticated';

var SitesRoute = AuthenticatedRoute.extend({
  model: function() {
    return this.store.findAll('site');
  }
});

export default SitesRoute;

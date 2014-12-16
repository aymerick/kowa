import AuthenticatedRoute from 'kowa/routes/authenticated';

var SiteRoute = AuthenticatedRoute.extend({
  model: function(params) {
    return this.store.find('site', params.site_id);
  }
});

export default SiteRoute;

import AuthenticatedRoute from 'kowa/routes/authenticated';

var SiteIndexRoute = AuthenticatedRoute.extend({
  afterModel: function () {
    // redirects /:site_id/ to /:site_id/posts
    this.transitionTo('posts');
  }
});

export default SiteIndexRoute;

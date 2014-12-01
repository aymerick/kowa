import AuthenticatedRoute from 'kowa/routes/authenticated';

var SiteRoute = AuthenticatedRoute.extend({
  afterModel: function (model, transition) {
    // redirects /:site_id to /:site_id/posts
    if (transition.targetName == "site.index") {
      this.transitionTo('posts');
    }
  }
});

export default SiteRoute;

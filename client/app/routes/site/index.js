import Ember from 'ember';
import AuthenticatedRoute from 'kowa/routes/authenticated';

var SiteIndexRoute = AuthenticatedRoute.extend({
  afterModel: function () {
    var site = this.modelFor('site');
    if (Ember.isBlank(site.get('description'))) {
      // redirects /:site_id/ to /:site_id/settings
      this.transitionTo('settings');
    } else {
      // redirects /:site_id/ to /:site_id/posts
      this.transitionTo('posts');
    }
  }
});

export default SiteIndexRoute;

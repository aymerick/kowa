import Ember from 'ember';
import AuthenticatedRoute from 'kowa/routes/authenticated';

var initialPaginationParams = { 'page': 1, 'perPage': 15 };

var PostsRoute = AuthenticatedRoute.extend({
  paginationParams: null,

  model: function() {
    var site = this.modelFor('site');
    var params = Ember.merge({ 'site': site.get('id') }, initialPaginationParams);

    this.set('paginationParams', params);

    // using filter() allows the template to auto-update when new models are pulled in from the server
    return this.store.filter('post', params, function () {
      // nothing to filter
      return true;
    });
  },

  afterModel: function (posts, transition) {
    // redirects /:site_id/posts to to /:site_id/posts/:post_id
    if (transition.targetName === "posts.index") {
      var firstPost = posts.get('firstObject');
      if (firstPost) {
        this.transitionTo('posts.post', firstPost);
      }
    }
  },

  setupController: function (controller, model) {
    this._super(controller, model);

    // setup pagination controller mixin
    controller.setupPagination('post', this.paginationParams);
  }
});

export default PostsRoute;

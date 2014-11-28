import AuthenticatedRoute from 'kowa/routes/authenticated';

var PostsRoute = AuthenticatedRoute.extend({
  model: function() {
    var site = this.modelFor('site');
    var selector = { 'site': site.get('id'), 'page': 1, 'per_page': 10 };

    // using filter() allows the template to auto-update when new models are pulled in from the server
    return this.store.filter('post', selector, function () {
      // nothing to filter
      return true;
    });
  }
});

export default PostsRoute;

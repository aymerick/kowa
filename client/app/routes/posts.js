import AuthenticatedRoute from 'kowa/routes/authenticated';

var PostsRoute = AuthenticatedRoute.extend({
  model: function() {
    var site = this.modelFor('site');
    var selector = { 'site': site.get('id'), 'page': 1 };

    // using filter() allows the template to auto-update when new models are pulled in from the server
    return this.store.filter('post', selector, function () {
      return true;
    });
  }
});

export default PostsRoute;

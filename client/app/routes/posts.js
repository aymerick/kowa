import AuthenticatedRoute from 'kowa/routes/authenticated';

var PostsRoute = AuthenticatedRoute.extend({
  model: function(params) {
    return this.modelFor('site').get('posts');
  }
});

export default PostsRoute;

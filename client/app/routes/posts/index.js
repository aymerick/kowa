import AuthenticatedRoute from 'kowa/routes/authenticated';

var PostsIndexRoute = AuthenticatedRoute.extend({
  // redirects /:site_id/posts to /:site_id/posts/:post_id
  beforeModel: function () {
    var posts = this.modelFor('posts');
    var firstPost = posts.get('firstObject');
    if (firstPost) {
      this.transitionToRoute('post', firstPost);
    }
  }
});

export default PostsIndexRoute;

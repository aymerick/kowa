import AuthenticatedRoute from 'kowa/routes/authenticated';

var PostEditRoute = AuthenticatedRoute.extend({
  model: function(params) {
    return this.store.find('post', params.post_id);
  }
});

export default PostEditRoute;

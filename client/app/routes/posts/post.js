import AuthenticatedRoute from 'kowa/routes/authenticated';
import ContentEditionRoute from 'kowa/mixins/content-edition-route';

var PostsPostRoute = AuthenticatedRoute.extend(ContentEditionRoute, {
  model: function(params) {
    return this.store.find('post', params.post_id);
  }
});

export default PostsPostRoute;

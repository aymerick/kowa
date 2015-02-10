import AuthenticatedRoute from 'kowa/routes/authenticated';
import ContentListRoute from 'kowa/mixins/content-list-route';

var PostsRoute = AuthenticatedRoute.extend(ContentListRoute, {
  contentModelType: 'post'
});

export default PostsRoute;

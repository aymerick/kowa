import AuthenticatedRoute from 'kowa/routes/authenticated';
import PostEditionRoute from 'kowa/mixins/post-edition-route';

var PostEditRoute = AuthenticatedRoute.extend(PostEditionRoute, {
});

export default PostEditRoute;

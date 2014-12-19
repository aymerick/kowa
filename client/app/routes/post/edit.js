import AuthenticatedRoute from 'kowa/routes/authenticated';
import ContentEditionRoute from 'kowa/mixins/content-edition-route';
import PostEditionRoute from 'kowa/mixins/post-edition-route';

var PostEditRoute = AuthenticatedRoute.extend(ContentEditionRoute, PostEditionRoute, {
});

export default PostEditRoute;

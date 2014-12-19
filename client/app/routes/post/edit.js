import AuthenticatedRoute from 'kowa/routes/authenticated';
import ContentEditionRoute from 'kowa/mixins/content-edition-route';

var PostEditRoute = AuthenticatedRoute.extend(ContentEditionRoute, {
});

export default PostEditRoute;

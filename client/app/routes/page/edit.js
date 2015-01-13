import AuthenticatedRoute from 'kowa/routes/authenticated';
import ContentEditionRoute from 'kowa/mixins/content-edition-route';

var PageEditRoute = AuthenticatedRoute.extend(ContentEditionRoute, {
});

export default PageEditRoute;

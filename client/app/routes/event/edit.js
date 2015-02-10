import AuthenticatedRoute from 'kowa/routes/authenticated';
import ContentEditionRoute from 'kowa/mixins/content-edition-route';

var EventEditRoute = AuthenticatedRoute.extend(ContentEditionRoute, {
});

export default EventEditRoute;

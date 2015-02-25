import AuthenticatedRoute from 'kowa/routes/authenticated';
import ContentEditionRoute from 'kowa/mixins/content-edition-route';

var PagesPageRoute = AuthenticatedRoute.extend(ContentEditionRoute, {
  model: function(params) {
    return this.store.find('page', params.page_id);
  }
});

export default PagesPageRoute;

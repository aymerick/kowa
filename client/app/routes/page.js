import AuthenticatedRoute from 'kowa/routes/authenticated';

var PageRoute = AuthenticatedRoute.extend({
  model: function(params) {
    return this.store.find('page', params.page_id);
  }
});

export default PageRoute;

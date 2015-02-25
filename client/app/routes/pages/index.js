import AuthenticatedRoute from 'kowa/routes/authenticated';

var PagesIndexRoute = AuthenticatedRoute.extend({
  // redirects /:site_id/pages to /:site_id/pages/:page_id
  beforeModel: function () {
    var pages = this.modelFor('pages');
    var firstPage = pages.get('firstObject');
    if (firstPage) {
      this.transitionTo('pages.page', firstPage);
    }
  }
});

export default PagesIndexRoute;

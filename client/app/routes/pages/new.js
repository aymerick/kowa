import AuthenticatedRoute from 'kowa/routes/authenticated';
import ContentEditionRoute from 'kowa/mixins/content-edition-route';

var PagesNewRoute = AuthenticatedRoute.extend(ContentEditionRoute, {
  // use PagesPageController
  controllerName: 'pages.page',

  // this is a fresh new model
  model: function() {
    return this.store.createRecord('page', { site: this.modelFor('site') });
  },

  // use page template
  renderTemplate: function() {
    this.render('pages/page');
  }
});

export default PagesNewRoute;

import AuthenticatedRoute from 'kowa/routes/authenticated';
import ContentEditionRoute from 'kowa/mixins/content-edition-route';
import Page from 'kowa/models/page';

var PagesNewRoute = AuthenticatedRoute.extend(ContentEditionRoute, {
  // use PageEditController
  controllerName: 'page.edit',

  // this is a fresh new model
  model: function() {
    return this.store.createRecord('page', Page.newRecordAttrs({
      site: this.modelFor('site')
    }));
  },

  // use 'page/edit' template
  renderTemplate: function() {
    this.render('page/edit');
  }
});

export default PagesNewRoute;

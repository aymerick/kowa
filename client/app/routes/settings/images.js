import Ember from 'ember';
import AuthenticatedRoute from 'kowa/routes/authenticated';

var initialPaginationParams = { 'page': 1, 'perPage': 20 };

var SettingsImagesRoute = AuthenticatedRoute.extend({
  paginationParams: null,

  model: function() {
    var site = this.modelFor('site');
    var params = Ember.merge({ 'site': site.get('id') }, initialPaginationParams);

    this.set('paginationParams', params);

    // 1. Caching resultPromise permits to avoid sending a request to server everytime we transition into this route
    // 2. Using filter() allows the template to auto-update when new models are pulled in from the server
    var resultPromise = this.get('resultPromise') || this.store.filter('image', params, function (image) {
      return image.get('site.id') === site.get('id');
    });

    this.set('resultPromise', resultPromise);

    return resultPromise;
  },

  setupController: function (controller, model) {
    this._super(controller, model);

    // setup pagination controller mixin
    controller.setupPagination('image', this.paginationParams);
  }
});

export default SettingsImagesRoute;

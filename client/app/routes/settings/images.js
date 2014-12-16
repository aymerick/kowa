import Ember from 'ember';
import AuthenticatedRoute from 'kowa/routes/authenticated';

var initialPaginationParams = { 'page': 1, 'perPage': 20 };

var SettingsImagesRoute = AuthenticatedRoute.extend({
  paginationParams: null,

  model: function() {
    var site = this.modelFor('site');
    var params = Ember.merge({ 'site': site.get('id') }, initialPaginationParams);

    this.set('paginationParams', params);

    // using filter() allows the template to auto-update when new models are pulled in from the server
    return this.store.filter('image', params, function () {
      // nothing to filter
      return true;
    });
  },

  setupController: function (controller, model) {
    this._super(controller, model);

    // setup pagination controller mixin
    controller.setupPagination('image', this.paginationParams);
  }
});

export default SettingsImagesRoute;

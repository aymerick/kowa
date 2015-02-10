import Ember from 'ember';

var initialPaginationParams = { 'page': 1, 'perPage': 15 };

var ContentListRouteMixin = Ember.Mixin.create({
  paginationParams: null,
  contentModelType: null,

  model: function() {
    var site = this.modelFor('site');
    var params = Ember.merge({ 'site': site.get('id') }, initialPaginationParams);

    this.set('paginationParams', params);

    if (this.get('contentModelType') === null) {
      Ember.assert('Content model type must be set.');
    }

    // 1. Caching resultPromise permits to avoid sending a request to server everytime we transition into this route
    // 2. Using filter() allows the template to auto-update when new models are pulled in from the server
    var resultPromise = this.get('resultPromise') || this.store.filter(this.get('contentModelType'), params, function () {
      // nothing to filter
      return true;
    });

    this.set('resultPromise', resultPromise);

    return resultPromise;
  },

  setupController: function (controller, model) {
    this._super(controller, model);

    // setup pagination controller mixin
    controller.setupPagination(this.get('contentModelType', this.paginationParams));
  }
});

export default ContentListRouteMixin;

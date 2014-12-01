import Ember from 'ember';

var PaginationControllerMixin = Ember.Mixin.create({
  modelType: null,
  paginationParams: null,

  nextPage: null,
  isLoading: false,

  // setup pagination
  setupPagination(modelType, paginationParams) {
    this.set('modelType', modelType);
    this.set('paginationParams', paginationParams);

    var metadata = this.store.metadataFor(modelType);
    this.set('nextPage', metadata.nextPage);
  },

  reportLoadError: function (/* response */) {
    // @todo !!!
  },

  actions: {
    loadNextPage: function () {
      var self = this;

      var nextPage = this.get('nextPage');
      if (nextPage) {
        this.set('isLoading', true);
        this.set('paginationParams.page', nextPage);

        this.store.find(this.get('modelType'), this.get('paginationParams')).then(function () {
          var metadata = self.store.metadataFor(self.get('modelType'));

          self.set('nextPage', metadata.nextPage);
          self.set('isLoading', false);
        }, function (response) {
          self.reportLoadError(response);
        });
      }
    }
  }
});

export default PaginationControllerMixin;

import Ember from 'ember';
import PaginationControllerMixin from 'kowa/mixins/pagination-controller';

var SelectImageModal = Ember.ArrayController.extend(PaginationControllerMixin, {
  sortProperties: ['createdAt'],
  sortAscending: false,

  selectedImage: null,

  nothingSelected: Ember.computed.not('selectedImage'),

  setupModal: function(masterController, field) {
    this.set('masterController', masterController);
    this.set('field', field);

    var site = masterController.get('site');

    var params = { 'site': site.get('id'), 'page': 1, 'perPage': 16 };
    this.setupPagination('image', params);

    this.set('model', this.store.filter('image', params, function () {
      // nothing to filter
      return true;
    }));
  },

  actions: {
    selectImage: function (image) {
      this.set('selectedImage', image);
    },

    save: function() {
      var masterController = this.get('masterController');
      var field = this.get('field');

      masterController.send('imageSelected', field, this.get('selectedImage'));
    }
  }
});

export default SelectImageModal;

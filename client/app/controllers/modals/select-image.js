import Ember from 'ember';
import PaginationControllerMixin from 'kowa/mixins/pagination-controller';

var SelectImageModal = Ember.ArrayController.extend(PaginationControllerMixin, {
  sortProperties: ['createdAt'],
  sortAscending: false,

  selectedImage: null,

  nothingSelected: Ember.computed.not('selectedImage'),

  masterController: null,
  masterField: null,

  setupModal: function(masterController, masterField) {
    var site = masterController.get('site');

    var params = { 'site': site.get('id'), 'page': 1, 'perPage': 16 };
    this.setupPagination('image', params);

    this.set('model', this.store.filter('image', params, function () {
      // nothing to filter
      return true;
    }));

    this.set('masterController', masterController);
    this.set('masterField', masterField);
  },

  actions: {
    selectImage: function (image) {
      this.set('selectedImage', image);
    },

    save: function() {
      var masterController = this.get('masterController');
      var masterField = this.get('masterField');

      masterController.send('imageSelected', masterField, this.get('selectedImage'));
    }
  },

  // i18n for attributes values and components parameters
  i18n: function() {
    return {
      chooseImage: this.t('chooseImage'),
      select: this.t('select'),
      cancel: this.t('cancel')
    };
  }.property() // @todo Watch current language
});

export default SelectImageModal;

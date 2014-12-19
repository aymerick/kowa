import Ember from 'ember';
import PaginationControllerMixin from 'kowa/mixins/pagination-controller';

var SelectImageModal = Ember.ArrayController.extend(PaginationControllerMixin, {
  sortProperties: ['createdAt'],
  sortAscending: false,

  selectedImage: null,
  targetModel: null,
  targetField: null,

  nothingSelected: Ember.computed.not('selectedImage'),

  setupModal: function(model, type) {
    var params = { 'site': model.get('id'), 'page': 1, 'perPage': 16 };

    var self = this;
    this.set('model', this.store.filter('image', params, function () {
      self.set('targetModel', model);
      self.set('targetField', type);
      self.setupPagination('image', params);

      // nothing to filter
      return true;
    }));
  },

  actions: {
    selectImage: function (image) {
      this.set('selectedImage', image);
    },

    save: function() {
      var self = this;

      // @todo FIXME ! This is really ugly !
      if (this.selectedImage && this.targetModel && this.targetField) {
        // update model
        this.targetModel.set(this.targetField, this.selectedImage);

        // save model
        this.targetModel.save().then(function (siteSaved) {
            self.get('flashes').success('Saved.');
            return siteSaved;
        }).catch(function () {
            self.get('flashes').danger('Failed to save.');
        });
      }
      else {
        self.get('flashes').danger('No image selected.');
      }
    }
  }
});

export default SelectImageModal;

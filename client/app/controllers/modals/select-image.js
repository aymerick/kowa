import Ember from 'ember';
import PaginationControllerMixin from 'kowa/mixins/pagination-controller';

var SelectImageModal = Ember.ArrayController.extend(PaginationControllerMixin, {
  sortProperties: ['createdAt'],
  sortAscending: false,

  actions: {
    save: function() {
      var self = this;

      this.get('model').save().then(function (model) {
          self.get('flashes').success('Saved.');
          return model;
      }).catch(function () {
          self.get('flashes').danger('Failed.');
      });
    }
  }
});

export default SelectImageModal;

import Ember from 'ember';

var DeleteImageModal = Ember.ObjectController.extend({
  actions: {
    delete: function() {
      var self = this;
      var model = this.get('model');

      model.destroyRecord().then(function () {
          self.get('flashes').success('Image deleted.');
      }).catch(function () {
          self.get('flashes').danger('Failed to delete image.');
          model.rollback();
          model.reload();
      });
    }
  }
});

export default DeleteImageModal;

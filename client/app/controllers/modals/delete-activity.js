import Ember from 'ember';

var DeleteActivityModal = Ember.ObjectController.extend({
  setupModal: function(masterController, activity) {
    this.set('model', activity);
  },

  actions: {
    delete: function() {
      var self = this;
      var model = this.get('model');

      model.destroyRecord().then(function () {
          self.get('flashes').success('Activity deleted.');
      }).catch(function () {
          self.get('flashes').danger('Failed to delete activity.');
      });
    }
  }
});

export default DeleteActivityModal;

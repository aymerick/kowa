import Ember from 'ember';

var LeaveEditionModal = Ember.ObjectController.extend({
  setupModal: function(model) {
    this.set('model', model);
  },

  actions: {
    confirm: function() {
      var args = this.get('model');
      var editionController, transition, model;

      // get arguments
      if (Ember.isArray(args)) {
          editionController = args[0];
          transition = args[1];
          model = editionController.get('model');
      }

      if (Ember.isNone(transition) || Ember.isNone(editionController) || Ember.isNone(model)) {
          this.get('flashes').danger('An error occured.');
          return true;
      }

      if (model.get('isNew')) {
          model.deleteRecord();
      } else {
          model.rollback();
      }

      // reset edition
      editionController.setupEdition();

      transition.retry();
    }
  }
});

export default LeaveEditionModal;

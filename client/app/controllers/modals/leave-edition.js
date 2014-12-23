import Ember from 'ember';

var LeaveEditionModal = Ember.ObjectController.extend({
  setupModal: function(masterController, haltedTransition) {
    this.set('masterController', masterController);
    this.set('haltedTransition', haltedTransition);
  },

  actions: {
    confirm: function() {
      var editionController, haltedTransition, model;

      // get arguments
      if (Ember.isArray(args)) {
          editionController = this.get('masterController');
          haltedTransition = this.get('transition');
          model = editionController.get('model');
      }

      if (Ember.isNone(haltedTransition) || Ember.isNone(editionController) || Ember.isNone(model)) {
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

      haltedTransition.retry();
    }
  }
});

export default LeaveEditionModal;

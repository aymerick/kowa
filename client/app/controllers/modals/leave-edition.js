import Ember from 'ember';

var LeaveEditionModal = Ember.ObjectController.extend({
  masterController: null,
  haltedTransition: null,

  setupModal: function(masterController, haltedTransition) {
    this.set('model', masterController.get('model'));
    this.set('masterController', masterController);
    this.set('haltedTransition', haltedTransition);
  },

  actions: {
    confirm: function() {
      var model = this.get('model');
      var masterController = this.get('masterController');
      var haltedTransition = this.get('haltedTransition');

      if (Ember.isNone(haltedTransition) || Ember.isNone(masterController) || Ember.isNone(model)) {
          this.get('flashes').danger('An error occured.');

          // bubble the action
          return true;
      }

      masterController.rollbackEdition();

      // this is necessary so that transition retry works
      masterController.set('isDirty', false);

      haltedTransition.retry();
    }
  },

  // i18n for attributes values and components parameters
  i18n: function() {
    return {
      leaveQuestion: this.t('edition.leaveQuestion'),
      leave: this.t('edition.leave'),
      stay: this.t('edition.stay')
    };
  }.property() // @todo Watch current language
});

export default LeaveEditionModal;

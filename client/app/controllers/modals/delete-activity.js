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
          self.get('flashes').success(self.t('activity.deleted'));
      }).catch(function () {
          self.get('flashes').danger(self.t('activity.deleteFailed'));
      });
    }
  },

  // i18n for attributes values and components parameters
  i18n: function() {
    return {
      deleteQuestion: this.t('activity.deleteQuestion'),
      'delete': this.t('delete'),
      cancel: this.t('cancel')
    };
  }.property() // @todo Watch current language
});

export default DeleteActivityModal;

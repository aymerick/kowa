import Ember from 'ember';

var DeleteMemberModal = Ember.ObjectController.extend({
  setupModal: function(masterController, member) {
    this.set('model', member);
  },

  actions: {
    delete: function() {
      var self = this;
      var model = this.get('model');

      model.destroyRecord().then(function () {
          self.get('flashes').success(self.t('member.deleted'));
      }).catch(function () {
          self.get('flashes').danger(self.t('member.deleteFailed'));
      });
    }
  },

  // i18n for attributes values and components parameters
  i18n: function() {
    return {
      deleteQuestion: this.t('member.deleteQuestion'),
      'delete': this.t('delete'),
      cancel: this.t('cancel')
    };
  }.property() // @todo Watch current language
});

export default DeleteMemberModal;

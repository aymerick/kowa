import Ember from 'ember';

var DeleteImageModal = Ember.ObjectController.extend({
  masterController: null,

  setupModal: function(masterController, image) {
    this.set('model', image);
    this.set('masterController', masterController);
  },

  actions: {
    delete: function() {
      var self = this;
      var model = this.get('model');

      model.destroyRecord().then(function () {
          self.get('flashes').success(self.t('image.deleted'));
      }).catch(function () {
          self.get('flashes').danger(self.t('image.deleteFailed'));
          model.rollback();
          model.reload();
      });
    }
  },

  // i18n for attributes values and components parameters
  i18n: function() {
    return {
      deleteQuestion: this.t('image.deleteQuestion'),
      'delete': this.t('delete'),
      cancel: this.t('cancel')
    };
  }.property() // @todo Watch current language
});

export default DeleteImageModal;

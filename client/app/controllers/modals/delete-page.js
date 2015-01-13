import Ember from 'ember';

var DeletePageModal = Ember.ObjectController.extend({
  nextRoute: null,

  setupModal: function(masterController, nextRoute) {
    this.set('model', masterController.get('model'));
    this.set('nextRoute', nextRoute);
  },

  actions: {
    delete: function() {
      var self = this;
      var model = this.get('model');

      model.destroyRecord().then(function () {
          self.get('flashes').success('Page deleted.');

          var nextRoute = self.get('nextRoute');
          if (!Ember.isNone(nextRoute)) {
            self.transitionToRoute(nextRoute);
          }
      }).catch(function () {
          self.get('flashes').danger('Failed to delete page.');
          model.rollback();
          model.reload();
      });
    }
  }
});

export default DeletePageModal;

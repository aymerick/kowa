import Ember from 'ember';

var DeletePostModal = Ember.ObjectController.extend({
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
          self.get('flashes').success('Post deleted.');

          var nextRoute = self.get('nextRoute');
          if (!Ember.isNone(nextRoute)) {
            self.transitionToRoute(nextRoute);
          }
      }).catch(function () {
          self.get('flashes').danger('Failed to delete post.');
          model.rollback();
          model.reload();
      });
    }
  }
});

export default DeletePostModal;

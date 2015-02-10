import Ember from 'ember';

var ContentDeleteModalMixin = Ember.Mixin.create({
  nextRoute: null,
  deleteMsgOk: null,
  deleteMsgFail: null,

  setupModal: function(masterController, nextRoute) {
    this.set('model', masterController.get('model'));
    this.set('nextRoute', nextRoute);
  },

  actions: {
    delete: function() {
      var self = this;
      var model = this.get('model');

      model.destroyRecord().then(function () {
          self.get('flashes').success(self.get('deleteMsgOk'));

          var nextRoute = self.get('nextRoute');
          if (!Ember.isNone(nextRoute)) {
            self.transitionToRoute(nextRoute);
          }
      }).catch(function () {
          self.get('flashes').danger(self.get('deleteMsgFail'));
          model.rollback();
          model.reload();
      });
    }
  }
});

export default ContentDeleteModalMixin;

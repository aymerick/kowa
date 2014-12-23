import Ember from 'ember';
import SimpleAuthApplicationRouteMixin from 'simple-auth/mixins/application-route-mixin';

var ApplicationRoute = Ember.Route.extend(SimpleAuthApplicationRouteMixin, {
  actions: {
    sessionAuthenticationFailed: function(/* error */) {
      this.get('controller').get('flashes').danger('Authentication failed.');
    },

    openModal: function(name, masterController, arg) {
      name = 'modals/' + name;

      var controller = this.controllerFor(name, true);
      if (controller && (Ember.typeOf(controller.setupModal) === 'function')) {
        controller.setupModal(masterController, arg);
      }

      this.render(name, {
        into: 'application',
        outlet: 'modal'
      });
    },

    closeModal: function() {
      this.disconnectOutlet({
        outlet: 'modal',
        parentView: 'application'
      });
    }
  }
});

export default ApplicationRoute;

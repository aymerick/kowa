import Ember from 'ember';
import SimpleAuthApplicationRouteMixin from 'simple-auth/mixins/application-route-mixin';

var ApplicationRoute = Ember.Route.extend(SimpleAuthApplicationRouteMixin, {
  actions: {
    sessionAuthenticationFailed: function(/* error */) {
      this.get('controller').get('flashes').danger('Authentication failed.');
    },

    openModal: function(name, model, type) {
      name = 'modals/' + name;

      if (this.controllerFor(name, true)) {
          this.controllerFor(name).set('model', model);
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

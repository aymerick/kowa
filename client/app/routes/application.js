import Ember from 'ember';
import SimpleAuthApplicationRouteMixin from 'simple-auth/mixins/application-route-mixin';

var ApplicationRoute = Ember.Route.extend(SimpleAuthApplicationRouteMixin, {
  actions: {
    sessionAuthenticationFailed: function(/* error */) {
      this.get('controller').get('flashes').danger('Authentication failed.');
    },

    openModal: function(name, model, type) {
      name = 'modals/' + name;

      var controller = this.controllerFor(name, true);
      if (controller) {
        // @todo FIXME ! This is really ugly !
        if (name === 'modals/select-image') {
          var params = { 'site': model.get('id'), 'page': 1, 'perPage': 16 };

          controller.set('model', this.store.filter('image', params, function () {
            controller.set('targetModel', model);
            controller.set('targetField', type);
            controller.setupPagination('image', params);

            // nothing to filter
            return true;
          }));
        }
        else {
          controller.set('model', model);
        }
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

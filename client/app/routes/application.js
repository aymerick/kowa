import Ember from 'ember';
import SimpleAuthApplicationRouteMixin from 'simple-auth/mixins/application-route-mixin';

var ApplicationRoute = Ember.Route.extend(SimpleAuthApplicationRouteMixin, {
  actions: {
    sessionAuthenticationFailed: function (error) {
      // this.controllerFor('application').set('loginErrorMessage', error.message);

      // if (error.errors) {
      //     this.notifications.showErrors(error.errors);
      // } else {
      //     // connection errors don't return proper status message, only req.body
      //     this.notifications.showError('There was a problem on the server.');
      // }
    }
  }
});

export default ApplicationRoute;

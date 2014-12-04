import Ember from 'ember';
import SimpleAuthApplicationRouteMixin from 'simple-auth/mixins/application-route-mixin';

var ApplicationRoute = Ember.Route.extend(SimpleAuthApplicationRouteMixin, {
  actions: {
    sessionAuthenticationFailed: function(/* error */) {
      this.get('controller').get('flashes').danger('Authentication failed.');
    }
  }
});

export default ApplicationRoute;

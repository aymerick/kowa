import Ember from 'ember';
import SimpleAuthUnauthenticatedRouteMixin from 'simple-auth/mixins/unauthenticated-route-mixin';

var UnauthenticatedRoute = Ember.Route.extend(SimpleAuthUnauthenticatedRouteMixin);

export default UnauthenticatedRoute;

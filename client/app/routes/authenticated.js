import Ember from 'ember';
import SimpleAuthAuthenticatedRouteMixin from 'simple-auth/mixins/authenticated-route-mixin';

var AuthenticatedRoute = Ember.Route.extend(SimpleAuthAuthenticatedRouteMixin);

export default AuthenticatedRoute;

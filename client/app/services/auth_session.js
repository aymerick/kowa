import Ember from 'ember';
import SimpleAuthSession from 'simple-auth/session';
import ajax from 'ic-ajax';

var AuthSession = SimpleAuthSession.extend({
  // adds currentUser property to simple-auth session
  currentUser: Ember.computed(function () {
    var self = this;

    var userPromise = ajax('/api/me').then(function(response) {
      var store = self.container.lookup('store:main');

      // push user to store
      var user = store.push('user', response.user);

      return user;
    });

    return Ember.ObjectProxy.extend(Ember.PromiseProxyMixin).create({
      promise: userPromise
    });
  })
});

export default AuthSession;

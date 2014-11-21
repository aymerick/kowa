import Ember from 'ember';
import SimpleAuthSession from 'simple-auth/session';

var AuthSession = SimpleAuthSession.extend({
  currentUser: Ember.computed(function () {
    return this.container.lookup('store:main').find('user', 'me');
  })
});

export default AuthSession;

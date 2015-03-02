import Ember from 'ember';
import LoginControllerMixin from 'simple-auth/mixins/login-controller-mixin';

var LoginController = Ember.Controller.extend(LoginControllerMixin, {
  authenticator: 'simple-auth-authenticator:oauth2-password-grant',

  // i18n for attributes values and components parameters
  i18n: function() {
    return {
      enterUsername: this.t('auth.enterUsername'),
      enterPassword: this.t('auth.enterPassword')
    };
  }.property() // @todo Watch current language
});

export default LoginController;

import ENV from '../config/environment';
import Ember from 'ember';
import AuthSession from 'kowa/services/auth_session';
import SimpleAuthOAuth2 from 'simple-auth-oauth2/authenticators/oauth2';

var AuthenticationInitializer = {
  name: 'authentication',
  before: 'simple-auth',

  initialize: function (container) {
    // custom session class
    container.register('session:kowa', AuthSession);

    // SimpleAuth
    ENV['simple-auth'] = {
      session: 'session:kowa',
      localStorageKey: 'kowa:session',
      authorizer: 'simple-auth-authorizer:oauth2-bearer',
      routeAfterAuthentication: 'sites',
      routeIfAlreadyAuthenticated: 'sites'
    };

    if (ENV.environment === 'development') {
      ENV['simple-auth']['crossOriginWhitelist'] = ['http://localhost:35830'];
    }

    ENV['simple-auth-oauth2'] = {
      serverTokenEndpoint: '/api/oauth/token',
      serverTokenRevocationEndpoint: '/api/oauth/revoke',
      refreshAccessTokens: true
    };

    SimpleAuthOAuth2.reopen({
      // Overrides method to inject oauth client id/secret in a basic auth header
      // ... this makes osin server happy
      makeRequest: function(url, data) {
        var client_id     = ENV.APP.api_client_id;
        var client_secret = ENV.APP.api_client_secret;

        return Ember.$.ajax({
          url:         url,
          type:        'POST',
          data:        data,
          dataType:    'json',
          contentType: 'application/x-www-form-urlencoded',
          headers:     { "Authorization": "Basic " + btoa(client_id + ":" + client_secret) }
        });
      }
    });
  }
};

export default AuthenticationInitializer;

import ENV from '../config/environment';
import AuthSession from 'kowa/services/auth_session';

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
  }
};

export default AuthenticationInitializer;

/* jshint node: true */

module.exports = function(environment) {
  var ENV = {
    modulePrefix: 'kowa',
    environment: environment,
    baseURL: '/',
    locationType: 'auto',
    EmberENV: {
      FEATURES: {
        // Here you can enable experimental features on an ember canary build
        // e.g. 'with-controller': true
      }
    },

    APP: {
      // Here you can pass flags/options to your application instance
      // when it is created
    }
  };

  // SimpleAuth
  ENV['simple-auth'] = {
    routeAfterAuthentication: 'sites',
    routeIfAlreadyAuthenticated: 'sites',
    localStorageKey: 'kowa:session',
    authorizer: 'simple-auth-authorizer:oauth2-bearer'
  }

  // SimpleAuth
  ENV['simple-auth-oauth2'] = {
    serverTokenEndpoint: '/oauth/token',
    // @todo serverTokenRevocationEndpoint: Ghost.apiRoot + '/oauth/revoke',
    refreshAccessTokens: true
  }

  if (environment === 'development') {
    // ENV.APP.LOG_RESOLVER = true;
    ENV.APP.LOG_ACTIVE_GENERATION = true;
    // ENV.APP.LOG_TRANSITIONS = true;
    // ENV.APP.LOG_TRANSITIONS_INTERNAL = true;
    ENV.APP.LOG_VIEW_LOOKUPS = true;

    // SimpleAuth
    ENV['simple-auth']['crossOriginWhitelist'] = ['http://localhost:35830']
  }

  if (environment === 'test') {
    // Testem prefers this...
    ENV.baseURL = '/';
    ENV.locationType = 'auto';

    // keep test console output quieter
    ENV.APP.LOG_ACTIVE_GENERATION = false;
    ENV.APP.LOG_VIEW_LOOKUPS = false;

    ENV.APP.rootElement = '#ember-testing';
  }

  if (environment === 'production') {

  }

  return ENV;
};

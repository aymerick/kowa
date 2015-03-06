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

      // thanks to ember-feature-flags addon
      FEATURES: {
        'show-content-format': false,
      },

      // oauth client id/secret (not used for now)
      api_client_id: 'kowa',
      api_client_secret: 'none',

      // ember-cli-i18n addon
      defaultLocale: 'en' // @todo i18n
    }
  };

  // cf. https://github.com/rwjblue/ember-cli-content-security-policy
  ENV.contentSecurityPolicy = {
    'default-src': "'none'",
    'script-src': "'self' 'unsafe-eval' 'unsafe-inline' localhost:35729 0.0.0.0:35729 127.0.0.1:35729",
    'font-src': "'self'",
    'connect-src': "'self' ws://localhost:35729 ws://0.0.0.0:35729 ws://127.0.0.1:35729",
    'img-src': "'self'",
    'style-src': "'self' 'unsafe-eval' 'unsafe-inline'",
    'media-src': "'self'"
  }

  if (environment === 'development') {
    // ENV.APP.LOG_RESOLVER = true;
    // ENV.APP.LOG_ACTIVE_GENERATION = true;
    // ENV.APP.LOG_TRANSITIONS = true;
    // ENV.APP.LOG_TRANSITIONS_INTERNAL = true;
    // ENV.APP.LOG_VIEW_LOOKUPS = true;
  }

  if (environment === 'test') {
    // Testem prefers this...
    ENV.baseURL = '/';
    ENV.locationType = 'auto';

    // keep test console output quieter
    ENV.APP.LOG_ACTIVE_GENERATION = false;
    ENV.APP.LOG_VIEW_LOOKUPS = false;

    ENV.APP.rootElement = '#ember-testing';

    ENV['simple-auth'] = {
      store: 'simple-auth-session-store:ephemeral'
    }
  }

  if (environment === 'production') {

  }

  return ENV;
};

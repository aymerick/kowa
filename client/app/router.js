import Ember from 'ember';
import config from './config/environment';

var Router = Ember.Router.extend({
  location: config.locationType
});

Router.map(function() {
  this.route('login');

  this.resource('sites', function() {
    this.route('site', { path: ':site_id' });
  });
});

export default Router;

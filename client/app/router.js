import Ember from 'ember';
import config from './config/environment';

var Router = Ember.Router.extend({
  location: config.locationType
});

Router.map(function() {
  this.route('login');
  this.route('logout');
  this.route('forgotten');

  this.resource('sites', {path: '/'}, function() {
    this.resource('site', { path: '/sites/:site_id' }, function() {
      this.resource('posts', {path: '/'}, function () {
        this.route('post', {path: ':post_id'});
      });

      this.resource('events', function () {
        this.route('event', {path: ':event_id'});
      });

      this.resource('pages', function () {
        this.route('page', {path: ':page_id'});
      });

      this.resource('settings', function () {
        this.route('general');

        this.resource('settings.actions', {path: '/actions'}, function () {
            this.route('action', {path: ':action_id'});
        });
      });
    });
  });
});

export default Router;

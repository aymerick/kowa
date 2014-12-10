import Ember from 'ember';
import config from './config/environment';

var Router = Ember.Router.extend({
  location: config.locationType
});

Router.map(function() {
  // auth
  this.route('login');
  this.route('logout');
  this.route('forgotten');

  this.resource('sites', {path: '/'}, function() {
    // site
    this.resource('site', {path: ':site_id'}, function() {
      // content
      this.resource('posts', function () {
        this.route('post', {path: ':post_id'});
      });

      this.resource('events', function () {
        this.route('event', {path: ':event_id'});
      });

      this.resource('pages', function () {
        this.route('page', {path: ':page_id'});
      });

      // content editor
      this.resource('editor', function() {
        this.route('post', {path: '/post'});
        this.route('post.edit', {path: '/post/:post_id'});

        this.route('event', {path: '/event'});
        this.route('event.edit', {path: '/event/:event_id'});

        this.route('page', {path: '/page'});
        this.route('page.edit', {path: '/page/:page_id'});
      });

      // settings
      this.resource('settings', function () {
        this.route('general');

        this.resource('settings.actions', {path: '/actions'}, function () {
            this.route('action', {path: ':action_id'});
        });

        this.resource('settings.members', {path: '/members'}, function () {
            this.route('member', {path: ':member_id'});
        });

        this.route('images');
      });
    });
  });

  // error
  this.route('error404', {path: '/*path'});
});

export default Router;

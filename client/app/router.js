import Ember from 'ember';
import config from './config/environment';

var Router = Ember.Router.extend({
  location: config.locationType
});

Router.map(function() {
  // auth
  this.route('login');
  this.route('logout'); // @todo
  this.route('forgotten'); // @todo

  // sites
  this.resource('sites', {path: '/'}, function() {
    // site
    this.resource('site', {path: ':site_id'}, function() {
      // content
      this.resource('posts', function () {
        this.route('new');
        this.route('post', {path: ':post_id'});
      });

      this.resource('events', function () {
        this.route('new');
        this.route('event', {path: ':event_id'});
      });

      this.resource('pages', function () {
        this.route('new');
        this.route('page', {path: ':page_id'});
      });

      // settings
      this.resource('settings', function () {
        this.route('general');

        this.route('homepage');

        this.resource('settings.activities', {path: '/activities'}, function () {
            this.route('page');

            this.route('new');
            this.route('activity', {path: ':activity_id'});
        });

        this.resource('settings.members', {path: '/members'}, function () {
            this.route('page');

            this.route('new');
            this.route('member', {path: ':member_id'});
        });

        this.resource('settings.contact', {path: '/contact'}, function () {
            this.route('page');
        });

        this.route('posts');

        this.route('events');

        this.route('images');

        this.route('account');
      });
    });
  });

  // error
  this.route('error404', {path: '/*path'});
});

export default Router;

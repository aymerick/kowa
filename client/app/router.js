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
        this.resource('post', {path: ':post_id'}, function() {
          this.route('edit', {path: '/edit'});
        });
      });

      // this.resource('events', function () {
      //   this.route('event', {path: ':event_id'});
      // });

      this.resource('pages', function () {
        this.route('new');
        this.resource('page', {path: ':page_id'}, function() {
          this.route('edit', {path: '/edit'});
        });
      });

      // settings
      this.resource('settings', function () {
        this.route('general');

        this.resource('settings.activities', {path: '/activities'}, function () {
            this.route('activity', {path: ':activity_id'});
        });

        // this.resource('settings.members', {path: '/members'}, function () {
        //   this.route('member', {path: ':member_id'});
        // });

        this.route('images');
      });
    });
  });

  // error
  this.route('error404', {path: '/*path'});
});

export default Router;

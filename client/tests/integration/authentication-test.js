import Ember from "ember";
import { test, moduleFor } from 'ember-qunit';
import startApp from '../helpers/start-app';
import Pretender from 'pretender';

var App;
var server;

moduleFor('route:sites', {
  beforeEach: function() {
    App = startApp();
    invalidateSession();

    var currentUserResp = {
      "user": {
        "id": "test",
        "createdAt": "2014-12-26T16:53:48.74+01:00",
        "email": "trucmush@wanadoo.fr",
        "firstName": "Jean-Claude",
        "lastName": "Trucmush",
        "links": {
          "sites": "/api/users/test/sites"
        }
      }
    };

    var currentUserSitesResp = {
      "sites": [{
        "id": "site_1",
        "createdAt": "2014-12-26T16:53:48.74+01:00",
        "updatedAt": "2015-02-06T18:13:49.995+01:00",
        "changedAt": "2015-02-09T10:51:40.568+01:00",
        "builtAt": "2015-02-09T10:51:40.661+01:00",
        "user": "test",
        "name": "My Site",
        "tagline": "This is my site bro",
        "description": "",
        "moreDesc": "",
        "joinText": "",
        "email": "",
        "address": "",
        "facebook": "e",
        "twitter": "",
        "googlePlus": "",
        "pageSettings": [],
        "theme": "willy",
        "uglyUrl": false,
        "links": {
          "activities": "/api/sites/site_1/activities",
          "events": "/api/sites/site_1/events",
          "images": "/api/sites/site_1/images",
          "pages": "/api/sites/site_1/pages",
          "posts": "/api/sites/site_1/posts"
        }
      }]
    };

    var site1Posts = {
      "meta": {
        "nextPage": null,
        "page": 1,
        "pages": 1,
        "perPage": 15,
        "prevPage": null,
        "skip": 0,
        "total": 1
      },
      "posts": [{
        "id": "54c66312c25c1908f7000042",
        "createdAt": "2015-01-25T16:53:48.74+01:00",
        "updatedAt": "2015-01-25T16:53:48.74+01:00",
        "site": "site_1",
        "publishedAt": "2015-01-25T16:53:48.74+01:00",
        "title": "This is the first post",
        "body": "Please, write something usefull"
      }]
    };

    server = new Pretender(function() {
      this.get("/api/me", function(request) {
        return [200, {"Content-Type": "application/json"}, JSON.stringify(currentUserResp)];
      });
      this.get("/api/users/test/sites", function(request) {
        return [200, {"Content-Type": "application/json"}, JSON.stringify(currentUserSitesResp)];
      });
      // ?site=site_1&page=1&perPage=15
      this.get("/api/posts", function(request) {
        return [200, {"Content-Type": "application/json"}, JSON.stringify(site1Posts)];
      });
    });
  },

  afterEach: function() {
    Ember.run(App, App.destroy);
    server.shutdown();
  }
});

test('Homepage redirects to login page if not authenticated yet', function(assert) {
  assert.expect(1);

  visit('/');

  andThen(function() {
    assert.equal(currentRouteName(), 'login');
  });
});

test('Homepage displays first post of first site if authenticated', function(assert) {
  assert.expect(1);

  authenticateSession();
  visit('/');

  andThen(function() {
    assert.equal(currentRouteName(), 'post.index');
  });
});

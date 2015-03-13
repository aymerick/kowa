import Ember from 'ember';
import AuthenticatedRoute from 'kowa/routes/authenticated';

var SitesRoute = AuthenticatedRoute.extend({
  model: function() {
    var self = this;

    return this.session.get('currentUser').then(function(currentUser){
      // setup locale
      var application = self.container.lookup('application:main');
      application.set('locale', currentUser.get('lang'));

      // observers locale change
      Ember.addObserver(currentUser, 'lang', function() {
        application.set('locale', currentUser.get('lang'));
      });

      return currentUser.get('sites');
    });
  },

  afterModel: function(sites, transition) {
    // redirects / to /:site_id
    if (!transition.params.site) {
      var firstSite = sites.get('lastObject');
      if (firstSite) {
        this.transitionTo('site', firstSite);
      }
    }
  }
});

export default SitesRoute;

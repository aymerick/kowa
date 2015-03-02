import Ember from 'ember';

var SiteController = Ember.ObjectController.extend({
  needs: ['sites'],
  userSites: Ember.computed.alias('controllers.sites.model'),
  selectedSite: Ember.computed.oneWay('model'),

  hideNav: Ember.computed.match('currentPath', /(error|login|forgotten)/),

  selectedSiteChanged: function() {
    if (this.get('selectedSite').get('id') !== this.get('model').get('id')) {
      // hard load to reset all controllers
      // @todo Find a better way to do that !
      window.location = window.location.origin + "/" + this.get('selectedSite').get('id');
      // this.transitionToRoute('site', this.get('selectedSite'));
      // this.unloadModels();
    }
  }.observes('selectedSite')
});

export default SiteController;

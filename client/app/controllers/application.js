import Ember from 'ember';

var ApplicationController = Ember.Controller.extend({
  // @todo Get that list from the server
  allLangs: [{
    id: 'en',
    name: 'English',
  }, {
    id: 'fr',
    name: 'Français',
  }]
});

export default ApplicationController;

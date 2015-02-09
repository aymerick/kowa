import Ember from 'ember';

var SettingsActivitiesController = Ember.ArrayController.extend({
  sortProperties: ['createdAt'],
  sortAscending: false
});

export default SettingsActivitiesController;

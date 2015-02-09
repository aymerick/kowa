import Ember from 'ember';

var SettingsActivitiesActivityController = Ember.ObjectController.extend({
  actions: {
    removeCover: function() {
      this.get('model').set('cover', null);
    },

    // called by 'select-image' modal controller
    imageSelected: function(field, image) {
      var model = this.get('model');
      model.set(field, image);
    },

    save: function () {
      var self = this;

      return this.get('model').save().then(function (model) {
        self.get('flashes').success('Activity saved.');

        return model;
      }).catch(function (/* errors */) {
        self.get('flashes').danger('Failed to save activity.');
      });
    }
  }
});

export default SettingsActivitiesActivityController;

import Ember from 'ember';

var SettingsActivitiesActivityController = Ember.ObjectController.extend({
  isSaving: false,

  actions: {
    removeCover: function() {
      this.get('model').set('cover', null);
    },

    // called by 'select-image' modal controller
    imageSelected: function(field, image) {
      var model = this.get('model');
      model.set(field, image);
    },

    // called by 'tinymce-editor' component
    summaryChanged: function(newValue) {
      this.get('model').set('summary', newValue);
    },

    // called by 'tinymce-editor' component
    bodyChanged: function(newValue) {
      this.get('model').set('body', newValue);
    },

    save: function () {
      var self = this;

      this.set('isSaving', true);

      return this.get('model').save().then(function (model) {
        self.get('flashes').success(self.t('activity.saved'));

        return model;
      }).catch(function (/* errors */) {
        self.get('flashes').danger(self.t('activity.saveFailed'));
      }).finally(function(){
        self.set('isSaving', false);
      });
    }
  }
});

export default SettingsActivitiesActivityController;

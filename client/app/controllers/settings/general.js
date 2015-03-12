import Ember from 'ember';

var SettingsGeneralController = Ember.ObjectController.extend({
  needs: ['settings'],
  site: Ember.computed.alias('controllers.settings.model'),
  isSaving: false,

  // @todo Get that list from the server
  allThemes: [ 'ailes', 'willy' ],

  // @todo Get that list from the server
  allLangs: [ 'en', 'fr' ],

  actions: {
    removeLogo: function() {
      this.get('model').set('logo', null);
    },

    // called by 'select-image' modal controller
    imageSelected: function(field, image) {
      var model = this.get('model');
      model.set(field, image);
    },

    save: function () {
      var self = this;

      this.set('isSaving', true);

      return this.get('model').save().then(function (model) {
        self.get('flashes').success('Settings saved.');

        return model;
      }).catch(function (/* errors */) {
        self.get('flashes').danger('Failed to save settings.');
      }).finally(function(){
        self.set('isSaving', false);
      });
    }
  }
});

export default SettingsGeneralController;

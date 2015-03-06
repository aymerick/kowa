import Ember from 'ember';

var SettingsGeneralController = Ember.ObjectController.extend({
  needs: ['settings'],
  site: Ember.computed.alias('controllers.settings.model'),
  isSaving: false,

  // @todo Get that list from the server
  allThemes: [ 'willy' ],

  actions: {
    removeLogo: function() {
      this.get('model').set('logo', null);
    },

    removeCover: function() {
      this.get('model').set('cover', null);
    },

    // called by 'select-image' modal controller
    imageSelected: function(field, image) {
      var model = this.get('model');
      model.set(field, image);
    },

    // called by 'tinymce-editor' component
    descriptionChanged: function(newValue) {
      this.get('model').set('description', newValue);
    },

    // called by 'tinymce-editor' component
    moreDescChanged: function(newValue) {
      this.get('model').set('moreDesc', newValue);
    },

    // called by 'tinymce-editor' component
    joinTextChanged: function(newValue) {
      this.get('model').set('joinText', newValue);
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

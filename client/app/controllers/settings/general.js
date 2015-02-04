import Ember from 'ember';

var SettingsGeneralController = Ember.ObjectController.extend({
    needs: ['settings'],
    site: Ember.computed.alias('controllers.settings.model'),

    // @todo Get that list from the server
    allThemes: [ 'test', 'willy' ],

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

        save: function () {
            var self = this;

            return this.get('model').save().then(function (model) {
                self.get('flashes').success('Settings saved.');

                return model;
            }).catch(function (/* errors */) {
                self.get('flashes').danger('Failed to save settings.');
            });
        }
    }
});

export default SettingsGeneralController;

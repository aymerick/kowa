import Ember from 'ember';

var SettingsGeneralController = Ember.ObjectController.extend({
    needs: ['settings'],
    site: Ember.computed.alias('controllers.settings.model'),

    actions: {
        removeLogo: function() {
            this.get('model').set('logo', null);
            this.send('save');
        },

        removeCover: function() {
            this.get('model').set('cover', null);
            this.send('save');
        },

        // called by 'select-image' modal controller
        imageSelected: function(field, image) {
            var model = this.get('model');
            model.set(field, image);

            // save model
            var self = this;
            model.save().then(function (siteSaved) {
                self.get('flashes').success('Saved.');
                return siteSaved;
            }).catch(function () {
                self.get('flashes').danger('Failed to save.');
            });
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

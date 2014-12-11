import Ember from 'ember';

var SettingsGeneralController = Ember.ObjectController.extend({
    actions: {
        removeLogo: function() {
            this.get('model').set('logo', null);
            this.send('save');
        },

        removeCover: function() {
            this.get('model').set('cover', null);
            this.send('save');
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

import Ember from 'ember';

var SettingsGeneralController = Ember.ObjectController.extend({
    actions: {
        save: function () {
            var self = this;

            return this.get('model').save().then(function (model) {
                self.get('flashes').success('Settings successfully saved.');

                return model;
            }).catch(function (/* errors */) {
                self.get('flashes').danger('Failed to save settings.');
            });
        }
    }
});

export default SettingsGeneralController;

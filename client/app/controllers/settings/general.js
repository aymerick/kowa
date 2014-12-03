import Ember from 'ember';

var SettingsGeneralController = Ember.ObjectController.extend({
    actions: {
        save: function () {
            var self = this;

            return this.get('model').save().then(function (model) {
                alert('saved');
                // @todo self.notifications.showSuccess('Settings successfully saved.');

                return model;
            }).catch(function (errors) {
                alert('save FAILED !');
                // @todo self.notifications.showErrors(errors);
            });
        }
    }
});

export default SettingsGeneralController;

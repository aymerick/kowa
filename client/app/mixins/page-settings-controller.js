import Ember from 'ember';

var PageSettingsControllerMixin = Ember.Mixin.create({
  needs: ['settings'],
  site: Ember.computed.alias('controllers.settings.model'),

  isSaving: false,

  pageSettingsSaveMsgOk: 'pageSettings.saved', // This is a i18n key
  pageSettingsSaveMsgErr: 'pageSettings.saveFailed', // This is a i18n key

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
      var model = this.get('model');
      var wasNew = model.get('isNew');

      this.set('isSaving', true);

      return model.save().then(function (savedModel) {
        if (wasNew) {
          // add to relationship
          self.get('site').get('pageSettings').addObject(savedModel);
        }

        self.get('flashes').success(self.t(self.get('pageSettingsSaveMsgOk')));

        return savedModel;
      }).catch(function (/* errors */) {
        self.get('flashes').danger(self.t(self.get('pageSettingsSaveMsgErr')));
      }).finally(function(){
        self.set('isSaving', false);
      });
    }
  }
});

export default PageSettingsControllerMixin;

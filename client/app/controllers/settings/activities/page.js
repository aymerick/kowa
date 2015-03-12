import Ember from 'ember';

var SettingsActivitiesPageController = Ember.ObjectController.extend({
  needs: ['settings'],
  site: Ember.computed.alias('controllers.settings.model'),

  emptyTitle: Ember.computed.empty('model.title'),
  cannotSave: Ember.computed.and('model.isNew', 'emptyTitle'),
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

        self.get('flashes').success(self.t('activity.pageSettingsSaved'));

        return savedModel;
      }).catch(function (/* errors */) {
        self.get('flashes').danger(self.t('activity.pageSettingsSaveFailed'));
      }).finally(function(){
        self.set('isSaving', false);
      });
    }
  }
});

export default SettingsActivitiesPageController;

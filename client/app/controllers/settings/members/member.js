import Ember from 'ember';

var SettingsMembersMemberController = Ember.ObjectController.extend({
  actions: {
    removePhoto: function() {
      this.get('model').set('photo', null);
    },

    // called by 'select-image' modal controller
    imageSelected: function(field, image) {
      var model = this.get('model');
      model.set(field, image);
    },

    save: function () {
      var self = this;

      return this.get('model').save().then(function (model) {
        self.get('flashes').success(self.t('member.saved'));

        return model;
      }).catch(function (/* errors */) {
        self.get('flashes').danger(self.t('member.saveFailed'));
      });
    }
  }
});

export default SettingsMembersMemberController;

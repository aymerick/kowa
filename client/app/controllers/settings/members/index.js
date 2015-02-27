import Ember from 'ember';

var SettingsMembersController = Ember.ArrayController.extend({
  sortProperties: ['fullname'],

  // i18n for attributes values and components parameters
  i18n: function() {
    return {
      photo: this.t('photo'),
      memberEdit: this.t('member.edit')
    };
  }.property() // @todo Watch current language
});

export default SettingsMembersController;

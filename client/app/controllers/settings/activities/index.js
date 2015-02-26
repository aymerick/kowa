import Ember from 'ember';

var SettingsActivitiesController = Ember.ArrayController.extend({
  sortProperties: ['createdAt'],
  sortAscending: false,

  // i18n for attributes values components parameters
  i18n: function() {
    return {
      coverImage: this.t('coverImage'),
      activityEdit: this.t('activityEdit')
    };
  }.property() // @todo Watch current language
});

export default SettingsActivitiesController;

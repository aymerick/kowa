import Ember from 'ember';
import PageSettingsControllerMixin from 'kowa/mixins/page-settings-controller';

var SettingsActivitiesPageController = Ember.ObjectController.extend(PageSettingsControllerMixin, {
  pageSettingsSaveMsgOk: 'activity.pageSettingsSaved', // This is a i18n key
  pageSettingsSaveMsgErr: 'activity.pageSettingsSaveFailed' // This is a i18n key
});

export default SettingsActivitiesPageController;

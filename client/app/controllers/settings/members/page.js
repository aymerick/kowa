import Ember from 'ember';
import PageSettingsControllerMixin from 'kowa/mixins/page-settings-controller';

var SettingsMembersPageController = Ember.ObjectController.extend(PageSettingsControllerMixin, {
  pageSettingsSaveMsgOk: 'member.pageSettingsSaved', // This is a i18n key
  pageSettingsSaveMsgErr: 'member.pageSettingsSaveFailed' // This is a i18n key
});

export default SettingsMembersPageController;

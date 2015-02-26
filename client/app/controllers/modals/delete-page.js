import Ember from 'ember';
import ContentDeleteModalMixin from 'kowa/mixins/content-delete-modal';

var DeletePageModal = Ember.ObjectController.extend(ContentDeleteModalMixin, {
  deleteMsgOk: 'page.delete', // This is an i18n key
  deleteMsgFail: 'page.deleteFailed', // This is an i18n key

  // i18n for attributes values and components parameters
  i18n: function() {
    return {
      deleteQuestion: this.t('page.deleteQuestion'),
      'delete': this.t('delete'),
      cancel: this.t('cancel')
    };
  }.property() // @todo Watch current language
});

export default DeletePageModal;

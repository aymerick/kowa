import Ember from 'ember';
import ContentDeleteModalMixin from 'kowa/mixins/content-delete-modal';

var DeleteEventModal = Ember.ObjectController.extend(ContentDeleteModalMixin, {
  deleteMsgOk: 'event.deleted', // This is an i18n key
  deleteMsgFail: 'event.deleteFailed', // This is an i18n key

  // i18n for attributes values and components parameters
  i18n: function() {
    return {
      deleteQuestion: this.t('event.deleteQuestion'),
      'delete': this.t('delete'),
      cancel: this.t('cancel')
    };
  }.property() // @todo Watch current language
});

export default DeleteEventModal;

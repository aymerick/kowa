import Ember from 'ember';
import ContentDeleteModalMixin from 'kowa/mixins/content-delete-modal';

var DeletePostModal = Ember.ObjectController.extend(ContentDeleteModalMixin, {
  deleteMsgOk: 'post.delete', // This is an i18n key
  deleteMsgFail: 'post.deleteFailed', // This is an i18n key

  // i18n for attributes values and components parameters
  i18n: function() {
    return {
      deleteQuestion: this.t('post.deleteQuestion'),
      'delete': this.t('delete'),
      cancel: this.t('cancel')
    };
  }.property() // @todo Watch current language
});

export default DeletePostModal;

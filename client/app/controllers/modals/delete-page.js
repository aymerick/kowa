import Ember from 'ember';
import ContentDeleteModalMixin from 'kowa/mixins/content-delete-modal';

var DeletePageModal = Ember.ObjectController.extend(ContentDeleteModalMixin, {
  deleteMsgOk: 'Page deleted.',
  deleteMsgFail: 'Failed to delete page.'
});

export default DeletePageModal;

import Ember from 'ember';
import ContentDeleteModalMixin from 'kowa/mixins/content-delete-modal';

var DeletePostModal = Ember.ObjectController.extend(ContentDeleteModalMixin, {
  deleteMsgOk: 'Post deleted.',
  deleteMsgFail: 'Failed to delete post.'
});

export default DeletePostModal;

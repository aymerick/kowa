import Ember from 'ember';
import ContentDeleteModalMixin from 'kowa/mixins/content-delete-modal';

var DeleteEventModal = Ember.ObjectController.extend(ContentDeleteModalMixin, {
  deleteMsgOk: 'Event deleted.',
  deleteMsgFail: 'Failed to delete event.'
});

export default DeleteEventModal;

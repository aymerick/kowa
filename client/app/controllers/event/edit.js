import Ember from 'ember';
import ContentEditionController from 'kowa/mixins/content-edition-controller';

var EventEditController = Ember.ObjectController.extend(ContentEditionController, {
  editionRelationships: Ember.A([ 'cover' ]),
  editionDefaultTitle: '(Untitled)',

  editionSaveMsgOk: 'Event saved.',
  editionSaveMsgErr: 'Failed to save event.'
});

export default EventEditController;

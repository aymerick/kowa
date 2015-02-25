import Ember from 'ember';
import ContentEditionController from 'kowa/mixins/content-edition-controller';

var EventsEventController = Ember.ObjectController.extend(ContentEditionController, {
  editionRelationships: Ember.A([ 'cover' ]),
  editionDefaultTitle: '(Untitled)',

  editionSaveMsgOk: 'Event saved.',
  editionSaveMsgErr: 'Failed to save event.',

  // @todo Get that list from the server
  allFormats: [
    { name: "Rich Text", id: 'html' },
    { name: "Markdown", id: 'md' }
  ]
});

export default EventsEventController;

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
  ],

  // i18n for attributes values components parameters
  i18n: function() {
    return {
      coverImage: this.t('coverImage'),
      eventDescription: this.t('eventDescription'),
      eventPlace: this.t('eventPlace'),
      eventTitle: this.t('eventTitle')
    };
  }.property() // @todo Watch current language
});

export default EventsEventController;

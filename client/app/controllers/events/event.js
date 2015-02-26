import Ember from 'ember';
import ContentEditionController from 'kowa/mixins/content-edition-controller';

var EventsEventController = Ember.ObjectController.extend(ContentEditionController, {
  editionRelationships: Ember.A([ 'cover' ]),
  editionDefaultTitle: 'event.untitled', // This is a i18n key

  editionSaveMsgOk: 'event.saved', // This is a i18n key
  editionSaveMsgErr: 'event.saveFailed', // This is a i18n key

  // @todo Get that list from the server
  allFormats: [
    { name: "Rich Text", id: 'html' },
    { name: "Markdown", id: 'md' }
  ],

  // i18n for attributes values components parameters
  i18n: function() {
    return {
      coverImage: this.t('coverImage'),
      eventDescription: this.t('event.description'),
      eventPlace: this.t('event.place'),
      eventTitle: this.t('event.title')
    };
  }.property() // @todo Watch current language
});

export default EventsEventController;

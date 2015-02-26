import Ember from 'ember';
import ContentEditionController from 'kowa/mixins/content-edition-controller';

var PagesPageController = Ember.ObjectController.extend(ContentEditionController, {
  editionRelationships: Ember.A([ 'cover' ]),
  editionDefaultTitle: 'page.untitled', // This is a i18n key

  editionSaveMsgOk: 'page.saved', // This is a i18n key
  editionSaveMsgErr: 'page.saveFailed', // This is a i18n key

  // @todo Get that list from the server
  allFormats: [
    { name: "Rich Text", id: 'html' },
    { name: "Markdown", id: 'md' }
  ],

  // i18n for attributes values components parameters
  i18n: function() {
    return {
      coverImage: this.t('coverImage'),
      pageBody: this.t('page.body'),
      pageTagline: this.t('page.tagline'),
      pageTitle: this.t('page.title')
    };
  }.property() // @todo Watch current language
});

export default PagesPageController;

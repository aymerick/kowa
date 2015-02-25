import Ember from 'ember';
import ContentEditionController from 'kowa/mixins/content-edition-controller';

var PagesPageController = Ember.ObjectController.extend(ContentEditionController, {
  editionRelationships: Ember.A([ 'cover' ]),
  editionDefaultTitle: '(Untitled)',

  editionSaveMsgOk: 'Page saved.',
  editionSaveMsgErr: 'Failed to save page.',

  // @todo Get that list from the server
  allFormats: [
    { name: "Rich Text", id: 'html' },
    { name: "Markdown", id: 'md' }
  ],

  // i18n for attributes values components parameters
  i18n: function() {
    return {
      coverImage: this.t('coverImage'),
      pageBody: this.t('pageBody'),
      pageTagline: this.t('pageTagline'),
      pageTitle: this.t('pageTitle')
    }
  }.property() // @todo Watch current language
});

export default PagesPageController;

import Ember from 'ember';
import ContentEditionController from 'kowa/mixins/content-edition-controller';

var PostsPostController = Ember.ObjectController.extend(ContentEditionController, {
  editionRelationships: Ember.A([ 'cover' ]),
  editionDefaultTitle: '(Untitled)',

  editionSaveMsgOk: 'Post saved.',
  editionSaveMsgErr: 'Failed to save post.',

  // @todo Get that list from the server
  allFormats: [
    { name: "Rich Text", id: 'html' },
    { name: "Markdown", id: 'md' }
  ],

  // i18n for attributes values components parameters
  i18n: function() {
    return {
      coverImage: this.t('coverImage'),
      postBody: this.t('postBody'),
      postTitle: this.t('postTitle')
    };
  }.property() // @todo Watch current language
});

export default PostsPostController;

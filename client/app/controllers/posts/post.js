import Ember from 'ember';
import ContentEditionController from 'kowa/mixins/content-edition-controller';

var PostsPostController = Ember.ObjectController.extend(ContentEditionController, {
  editionRelationships: Ember.A([ 'cover' ]),
  editionDefaultTitle: 'post.untitled', // This is a i18n key

  editionSaveMsgOk: 'post.saved', // This is a i18n key
  editionSaveMsgErr: 'post.saveFailed', // This is a i18n key

  // @todo Get that list from the server
  allFormats: [
    { name: "Rich Text", id: 'html' },
    { name: "Markdown", id: 'md' }
  ],

  // i18n for attributes values components parameters
  i18n: function() {
    return {
      coverImage: this.t('coverImage'),
      postBody: this.t('post.body'),
      postTitle: this.t('post.title')
    };
  }.property() // @todo Watch current language
});

export default PostsPostController;

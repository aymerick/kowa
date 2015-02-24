import Ember from 'ember';
import ContentEditionController from 'kowa/mixins/content-edition-controller';

var PostsPostController = Ember.ObjectController.extend(ContentEditionController, {
  editionRelationships: Ember.A([ 'cover' ]),
  editionDefaultTitle: '(Untitled)',

  editionSaveMsgOk: 'Post saved.',
  editionSaveMsgErr: 'Failed to save post.'
});

export default PostsPostController;

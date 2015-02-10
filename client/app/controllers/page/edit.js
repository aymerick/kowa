import Ember from 'ember';
import ContentEditionController from 'kowa/mixins/content-edition-controller';

var PageEditController = Ember.ObjectController.extend(ContentEditionController, {
  editionRelationships: Ember.A([ 'cover' ]),
  editionDefaultTitle: '(Untitled)',

  editionSaveMsgOk: 'Page saved.',
  editionSaveMsgErr: 'Failed to save page.',
});

export default PageEditController;

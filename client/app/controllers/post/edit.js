import Ember from 'ember';
import ContentEditionController from 'kowa/mixins/content-edition-controller';

var PostEditController = Ember.ObjectController.extend(ContentEditionController, {
  editionFields: Ember.A([ 'cover' ]),

  actions: {
    // called by 'select-image' modal controller
    imageSelected: function(field, image) {
      this.set(field, image);
    },

    removeCover: function() {
      this.get('model').set('cover', null);
    },

    savePost: function() {
      // set a default title
      var model = this.get('model');
      if (!model.get('title')) {
          model.set('title', '(Untitled)');
      }

      // save
      this.commitEdition('Post saved.', 'Failed to save post.');
    }
  }
});

export default PostEditController;

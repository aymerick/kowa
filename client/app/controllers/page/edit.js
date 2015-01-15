import Ember from 'ember';
import ContentEditionController from 'kowa/mixins/content-edition-controller';

var PageEditController = Ember.ObjectController.extend(ContentEditionController, {
  editionFields: Ember.A([ 'cover' ]),

  actions: {
    // called by 'select-image' modal controller
    imageSelected: function(field, image) {
      this.set(field, image);
    },

    removeCover: function() {
      this.get('model').set('cover', null);
    },

    savePage: function() {
      // set a default title
      var model = this.get('model');
      if (!model.get('title')) {
          model.set('title', '(Untitled)');
      }

      // save
      this.commitEdition('Page saved.', 'Failed to save page.');
    }
  }
});

export default PageEditController;

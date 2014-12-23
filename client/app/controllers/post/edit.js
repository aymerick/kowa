import Ember from 'ember';
import Post from 'kowa/models/post';

// @todo Refactor that into a ContentEditionController Mixin with properties:
//  - model fields that are editable (eg: ["title, "body"])
//  - msg to display on save success (eg: "Post saved.")
//  - msg to display on save error (eg: "Failed to save post.")
//  - new record properties: (eg: Post.newRecordAttrs())
var PostEditController = Ember.ObjectController.extend({
  isDirty: function() {
    var model = this.get('model');
    return ((model.get('isNew')) ||
            (model.get('title') !== this.get('titleScratch')) ||
            (model.get('body')  !== this.get('bodyScratch')) ||
            (model.get('cover') !== this.get('coverScratch')));
  }.property('titleScratch', 'bodyScratch', 'coverScratch', 'model.title', 'model.body', 'model.coverScratch', 'model.isNew'),

  nothingChanged: Ember.computed.not('isDirty'),

  setupEdition: function() {
    var model = this.get('model');

    this.set('titleScratch', model.get('title'));
    this.set('bodyScratch', model.get('body'));
    this.set('coverScratch', model.get('cover'));
  },

  actions: {
    // called by 'select-image' modal controller
    imageSelected: function(field, image) {
      this.set(field, image);
    },

    savePost: function() {
      if (!this.get('isDirty')) {
        // This should never happen
        return;
      }

      // update model
      var model = this.get('model');
      model.set('title', this.get('titleScratch'));
      model.set('body', this.get('bodyScratch'));
      model.set('cover', this.get('coverScratch'));

      // persist on server
      var self = this;
      model.save().then(function (postSaved) {
        self.get('flashes').success('Post saved.');

        self.transitionToRoute('post', postSaved);
        return postSaved;
      }).catch(function () {
        self.get('flashes').danger('Failed to save post.');

        // rollback model
        if (model.get('isNew')) {
          model.setProperties(Post.newRecordAttrs());
        }
        else {
          model.rollback();
        }
      });
    }
  }
});

export default PostEditController;

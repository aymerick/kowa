import Ember from 'ember';

var PostController = Ember.ObjectController.extend({
  isDirty: function() {
    return ((this.get('model').get('title') !== this.get('titleScratch')) ||
            (this.get('model').get('body')  !== this.get('bodyScratch')));
  }.property('titleScratch', 'bodyScratch', 'model.title', 'model.body'),

  actions: {
    editPost: function() {
      this.set('titleScratch', this.get('model').get('title') || "Untitled");
      this.set('bodyScratch', this.get('model').get('body'));

      this.set('isEditing', true);
    },

    savePost: function() {
      if (this.get('isDirty')) {
        // update model
        var model = this.get('model');
        model.set('title', this.get('titleScratch'));
        model.set('body', this.get('bodyScratch'));

        // persist on server
        var self = this;
        model.save().then(function (postSaved) {
            self.get('flashes').success('Post saved.');

            self.set('isEditing', false);

            return postSaved;
        }).catch(function () {
            self.get('flashes').danger('Failed to save post.');
        });
      }
      else {
        this.set('isEditing', false);
      }
    },

    cancelEdit: function() {
      this.set('isEditing', false);
    }
  }
});

export default PostController;

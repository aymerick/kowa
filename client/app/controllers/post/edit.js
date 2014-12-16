import Ember from 'ember';

var PostEditController = Ember.ObjectController.extend({
  isDirty: function() {
    return ((this.get('model').get('title') !== this.get('titleScratch')) ||
            (this.get('model').get('body')  !== this.get('bodyScratch')));
  }.property('titleScratch', 'bodyScratch', 'model.title', 'model.body'),

  actions: {
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

            self.transitionTo('post');
            return postSaved;
        }).catch(function () {
            self.get('flashes').danger('Failed to save post.');
        });
      }
      else {
        this.transitionTo('post');
      }
    },

    cancelEdit: function() {
      this.transitionTo('post');
    }
  }
});

export default PostEditController;

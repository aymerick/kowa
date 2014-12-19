import Ember from 'ember';
import Post from 'kowa/models/post';

var PostEditController = Ember.ObjectController.extend({
  isDirty: function() {
    var model = this.get('model');
    return ((model.get('isNew')) ||
            (model.get('title') !== this.get('titleScratch')) ||
            (model.get('body')  !== this.get('bodyScratch')));
  }.property('titleScratch', 'bodyScratch', 'model.title', 'model.body', 'model.isNew'),

  setupEdition: function() {
    var model = this.get('model');

    this.set('titleScratch', model.get('title'));
    this.set('bodyScratch', model.get('body'));
  },

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
      else {
        // @todo We should never reach that
        this.get('flashes').danger('Nothing to save.');
      }
    }
  }
});

export default PostEditController;

import Ember from 'ember';

var PostEditController = Ember.ObjectController.extend({
  isDirty: Ember.computed.alias('model.isDirty'),
  nothingChanged: Ember.computed.not('model.isDirty'),

  actions: {
    // called by 'select-image' modal controller
    imageSelected: function(field, image) {
      this.set(field, image);
    },

    removeCover: function() {
      this.get('model').set('cover', null);
    },

    savePost: function() {
      if (!this.get('isDirty')) {
        // This should never happen
        return;
      }

      // set a default title
      var model = this.get('model');
      if (!model.get('title')) {
          model.set('title', '(Untitled)');
      }

      // persist on server
      var self = this;
      model.save().then(function (postSaved) {
        self.get('flashes').success('Post saved.');

        return postSaved;
      }).catch(function () {
        self.get('flashes').danger('Failed to save post.');
      });
    }
  }
});

export default PostEditController;

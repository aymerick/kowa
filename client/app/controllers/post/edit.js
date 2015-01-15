import Ember from 'ember';

var PostEditController = Ember.ObjectController.extend({
  previousCover: null,

  isDirty: function() {
    var model = this.get('model');

    return (model.get('isDirty') || model.get('isNew') ||
            (model.get('cover') !== this.get('previousCover')));
  }.property('model.isDirty', 'model.isNew', 'model.cover', 'previousCover'),

  nothingChanged: Ember.computed.not('isDirty'),

  setupEdition: function() {
    var model = this.get('model');

    // cf. https://github.com/emberjs/data/issues/1308
    this.set('previousCover', model.get('cover'));
  },

  rollbackEdition: function() {
    var model = this.get('model');

    if (model.get('isNew')) {
        model.deleteRecord();
    } else {
        // cf. https://github.com/emberjs/data/issues/1308
        model.set('cover', this.get('previousCover'));

        model.rollback();
    }
  },

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

        self.set('previousCover', model.get('cover'));

        return postSaved;
      }).catch(function () {
        self.get('flashes').danger('Failed to save post.');
      });
    }
  }
});

export default PostEditController;

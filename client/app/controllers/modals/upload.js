import Ember from 'ember';

var UploadModal = Ember.ObjectController.extend({
  actions: {
    save: function() {
      var self = this;

      this.get('model').save().then(function (model) {
          self.get('flashes').success('Saved.');
          return model;
      }).catch(function (err) {
          self.get('flashes').danger('Failed.');
      });
    }
  }
});

export default UploadModal;

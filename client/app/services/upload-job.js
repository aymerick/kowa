import Ember from 'ember';

var UploadJob = Ember.Object.extend(Ember.Evented, {
  url: null,
  file: null,

  init: function(hash) {
    this._super(hash);
    this.setProperties(hash);
  }
});

export default UploadJob;

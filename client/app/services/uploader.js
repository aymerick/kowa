import Ember from "ember";
import EmberUploader from 'ember-uploader/uploader';

var Uploader = Ember.Object.extend({
  queue: Ember.A([ ]),
  currentJob: null,

  schedule: function(job) {
    this.get('queue').pushObject(job);

    Ember.run.next(this, this.runNextJob);
  },

  extractNextJob: function() {
    if (this.get('currentJob') !== null) {
      return null;
    }

    var job = this.get('queue').shiftObject();
    if (job !== null) {
      this.set('currentJob', job);
    }

    return job;
  },

  runNextJob: function() {
    var job = this.extractNextJob();
    if (job === null) {
      return;
    }

    var self = this;

    // create file uploader
    var fileUploader = EmberUploader.create({ url: job.url });

    fileUploader.on('progress', function(e) {
      job.trigger('uploadProgress', e.percent);
    });

    fileUploader.on('didUpload', function(response) {
      job.trigger('didUpload', response);

      // @todo on error !

      // cleanup
      delete(job.file);
      job.destroy();

      // run next job
      self.set('currentJob', null);
      Ember.run.next(self, self.runNextJob);
    });

    job.trigger('willUpload');

    fileUploader.upload(job.file);
  }
});

export default Uploader;

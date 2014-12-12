import Ember from 'ember';

import EmberUploader from 'ember-uploader/uploader';
import EmberUploaderFileField from 'ember-uploader/file-field';

var ImageUploader = EmberUploaderFileField.extend({
  multiple: true,
  url: '',
  accept: 'image/*',

  filesDidChange: (function() {
    var files = this.get('files');

    var uploader = EmberUploader.create({ url: this.get('url') });

    uploader.on('progress', function(e) {
      // Handle progress changes
      // Use `e.percent` to get percentage
    });

    uploader.on('didUpload', function(e) {
      // Handle finished upload
    });

    var filesLen = files.length;
    for (var i = 0; i < filesLen; i++) {
      uploader.upload(files[i]);
    }
  }).observes('files')
});

export default ImageUploader;

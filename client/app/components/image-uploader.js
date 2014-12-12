import EmberUploader from 'ember-uploader/uploader';
import EmberUploaderFileField from 'ember-uploader/file-field';

var ImageUploader = EmberUploaderFileField.extend({
  multiple: true,
  url: '',
  accept: 'image/*',

  filesDidChange: (function() {
    var files = this.get('files');

    for (var i = 0; i < files.length; i++) {
      this.uploadFile(files[i]);
    }
  }).observes('files'),

  uploadFile: function(file) {
    // lastModified: 1417793070000
    // lastModifiedDate: Fri Dec 05 2014 16:24:30 GMT+0100 (CET)
    // name: "Content.png"
    // size: 629704
    // type: "image/png"

    var uploader = EmberUploader.create({ url: this.get('url') });

    uploader.on('progress', function(e) {
      // loaded: 629888
      // percent: 100
      // position: 629888
      // total: 629888
      // totalSize: 629888
      // type: "progress"
    });

    uploader.on('didUpload', function(e) {
      // name: "Content.png"
      // path: "/upload/Content_rfBd56ti2SMtYvSg.png"
      // size: 629704
    });

    uploader.upload(file);
  }
});

export default ImageUploader;

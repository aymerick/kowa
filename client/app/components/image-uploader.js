import EmberUploaderFileField from 'ember-uploader/file-field';

var ImageUploader = EmberUploaderFileField.extend({
  multiple: true,
  accept: 'image/*',

  filesDidChange: (function() {
    var files = this.get('files');

    for (var i = 0; i < files.length; i++) {
      this.sendAction('upload', files[i], i);
    }
  }).observes('files')
});

export default ImageUploader;

import Ember from 'ember';

var ImageView = Ember.View.extend({
  tagName: 'li',
  classNames: ['thumbnail', 'images-grid-item'],
  classNameBindings: ['isSelected', 'content.isUploading'],

  isSelected: function() {
    var selectedImg = this.get('controller').get('selectedImage');
    return (selectedImg && selectedImg.get('id') === this.get('content').get('id'));
  }.property('controller.selectedImage'),

  click: function() {
    var image = this.get('content');
    if (!image.get('isUploading')) {
      this.get('controller').send('selectImage', image);
    }
  }
});

export default ImageView;

import Ember from 'ember';

var ImageView = Ember.View.extend({
  tagName: 'li',
  classNames: ['thumbnail', 'images-grid-item'],
  classNameBindings: ['isSelected'],

  isSelected: function() {
    var selectedImg = this.get('controller').get('selectedImage');
    return (selectedImg && selectedImg.get('id') === this.get('content').get('id'));
  }.property('controller.selectedImage'),

  click: function() {
    this.get('controller').send('selectImage', this.get('content'));
  }
});

export default ImageView;

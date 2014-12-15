import Ember from 'ember';

var SettingsImageView = Ember.View.extend({
  tagName: 'li',
  classNames: ['thumbnail', 'images-grid-item'],
  classNameBindings: 'content.isUploading'
});

export default SettingsImageView;

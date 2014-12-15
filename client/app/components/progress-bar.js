import Ember from "ember";

var ProgressBar = Ember.Component.extend({
  classNames: 'progress-bar progress-bar-success',
  attributeBindings: 'style',

  style: function() {
    return 'width: ' + this.get('progress') + '%;';
  }.property('progress')
});

export default ProgressBar;

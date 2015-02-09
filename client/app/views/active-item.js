import Ember from 'ember';

var ActiveItemView = Ember.View.extend({
  classNameBindings: ['active'],

  // set the 'active' class when our first child, that is a link, have the 'active' class set by the {{link-to}} ember magic
  active: Ember.computed('childViews.firstObject.active', function () {
    return this.get('childViews.firstObject.active');
  })
});

export default ActiveItemView;

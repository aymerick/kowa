// mixin used for routes that need to set a css className on the body tag
// source: https://github.com/TryGhost/Ghost/blob/master/core/client/mixins/style-body.js

import Ember from 'ember';

var StyleBodyMixin = Ember.Mixin.create({
  activate: function () {
    this._super();

    var cssClasses = this.get('classNames');

    if (cssClasses) {
      Ember.run.schedule('afterRender', null, function () {
        cssClasses.forEach(function (curClass) {
          Ember.$('body').addClass(curClass);
        });
      });
    }
  },

  deactivate: function () {
    this._super();

    var cssClasses = this.get('classNames');

    Ember.run.schedule('afterRender', null, function () {
      cssClasses.forEach(function (curClass) {
        Ember.$('body').removeClass(curClass);
      });
    });
  }
});

export default StyleBodyMixin;

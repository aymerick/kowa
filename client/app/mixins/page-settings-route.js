import Ember from 'ember';

var PageSettingsRouteMixin = Ember.Mixin.create({
  pageKind: null,

  model: function() {
    var site = this.modelFor('site');
    var self = this;

    var result = site.get('pageSettings').find(function(pageSettings) {
      return pageSettings.get('kind') === self.get('pageKind');
    });

    if (Ember.isNone(result)) {
      result = this.store.createRecord('site-page-setting', { kind: this.get('pageKind') });
    }

    // set fake attribute
    result.set('site', site);

    return result;
  },

  deactivate: function () {
    var model = this.modelFor('settings.' + this.get('pageKind') + '.page');

    // delete if not saved
    if (model && model.get('isNew')) {
        model.rollback();
    }

    this._super();
  }
});

export default PageSettingsRouteMixin;

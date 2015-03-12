import Ember from 'ember';
import AuthenticatedRoute from 'kowa/routes/authenticated';

var SettingsActivitiesPageRoute = AuthenticatedRoute.extend({
  model: function() {
    var site = this.modelFor('site');

    var result = site.get('pageSettings').find(function(pageSettings) {
      return pageSettings.get('kind') === 'activities';
    });

    if (Ember.isNone(result)) {
      result = this.store.createRecord('site-page-setting', { kind: 'activities' });
    }

    // set fake attribute
    result.set('site', site);

    return result;
  },

  deactivate: function () {
    var model = this.modelFor('settings.activities.page');

    // delete if not saved
    if (model && model.get('isNew')) {
        model.rollback();
    }

    this._super();
  }
});

export default SettingsActivitiesPageRoute;

import AuthenticatedRoute from 'kowa/routes/authenticated';

var SettingsActivitiesActivityRoute = AuthenticatedRoute.extend({
  model: function(params) {
    return this.store.find('activity', params.activity_id);
  },

  deactivate: function () {
    var model = this.modelFor('settings.activities.activity');

    // revert any unsaved changes on exit
    if (model && model.get('isDirty')) {
        model.rollback();
    }

    this._super();
  }
});

export default SettingsActivitiesActivityRoute;

import AuthenticatedRoute from 'kowa/routes/authenticated';
import Activity from 'kowa/models/activity';

var SettingsActivitiesNewRoute = AuthenticatedRoute.extend({
  // use PostEditController
  controllerName: 'settings.activities.activity',

  // this is a fresh new model
  model: function() {
    return this.store.createRecord('activity', Activity.newRecordAttrs({
      site: this.modelFor('site')
    }));
  },

  // use existing template
  renderTemplate: function() {
    this.render('settings/activities/activity');
  },

  deactivate: function () {
    var model = this.modelFor('settings.activities.new');

    // delete if not saved
    if (model && model.get('isNew')) {
        model.rollback();
    }

    this._super();
  }
});

export default SettingsActivitiesNewRoute;

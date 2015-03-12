import AuthenticatedRoute from 'kowa/routes/authenticated';

var SettingsContactIndexRoute = AuthenticatedRoute.extend({
  // use SettingsGeneralController
  controllerName: 'settings.general',

  model: function() {
    return this.modelFor('site');
  },

  deactivate: function () {
    var model = this.modelFor('settings.contact');

    // revert any unsaved changes on exit
    if (model && model.get('isDirty')) {
        model.rollback();
    }

    this._super();
  }
});

export default SettingsContactIndexRoute;

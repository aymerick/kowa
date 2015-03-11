import AuthenticatedRoute from 'kowa/routes/authenticated';

var SettingsHomepageRoute = AuthenticatedRoute.extend({
  model: function() {
    return this.modelFor('site');
  },

  deactivate: function () {
    var model = this.modelFor('settings.homepage');

    // revert any unsaved changes on exit
    if (model && model.get('isDirty')) {
        model.rollback();
    }

    this._super();
  }
});

export default SettingsHomepageRoute;

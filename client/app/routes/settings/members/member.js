import AuthenticatedRoute from 'kowa/routes/authenticated';

var SettingsMembersMemberRoute = AuthenticatedRoute.extend({
  model: function(params) {
    return this.store.find('member', params.member_id);
  },

  deactivate: function () {
    var model = this.modelFor('settings.members.member');

    // revert any unsaved changes on exit
    if (model && model.get('isDirty')) {
        model.rollback();
    }

    this._super();
  }
});

export default SettingsMembersMemberRoute;

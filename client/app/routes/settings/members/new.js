import AuthenticatedRoute from 'kowa/routes/authenticated';

var SettingsMembersNewRoute = AuthenticatedRoute.extend({
  // use SettingsMembersMemberRoute
  controllerName: 'settings.members.member',

  // this is a fresh new model
  model: function() {
    return this.store.createRecord('member', { site: this.modelFor('site') });
  },

  // use existing template
  renderTemplate: function() {
    this.render('settings/members/member');
  },

  deactivate: function () {
    var model = this.modelFor('settings.members.new');

    // delete if not saved
    if (model && model.get('isNew')) {
        model.rollback();
    }

    this._super();
  }
});

export default SettingsMembersNewRoute;

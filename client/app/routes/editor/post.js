import AuthenticatedRoute from 'kowa/routes/authenticated';

var EditorPostRoute = AuthenticatedRoute.extend({
  classNames: ['editor'],

  model: function () {
    var site = this.modelFor('site');

    return this.store.createRecord('post', {
      site: site
    });
  }
});

export default EditorPostRoute;

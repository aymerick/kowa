import AuthenticatedRoute from 'kowa/routes/authenticated';

var EditorIndexRoute = AuthenticatedRoute.extend({
  beforeModel: function () {
    this.transitionTo('editor.post');
  }
});

export default EditorIndexRoute;

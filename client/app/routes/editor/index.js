import Ember from 'ember';
import AuthenticatedRoute from 'kowa/routes/authenticated';

var EditorIndexRoute = Ember.Route.extend({
  beforeModel: function () {
    this.transitionTo('editor.post');
  }
});

export default EditorIndexRoute;

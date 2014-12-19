import AuthenticatedRoute from 'kowa/routes/authenticated';
import Post from 'kowa/models/post';

// Uses post/edit controller and template, but with a new record
var PostsNewRoute = AuthenticatedRoute.extend({
  controllerName: 'post.edit',

  model: function(params) {
    return this.store.createRecord('post', Post.newRecordAttrs({
      site: this.modelFor('site')
    }));
  },

  // @todo Factorize that with post/edit route
  setupController: function(controller, model) {
    this._super(controller, model);

    // init temp data
    controller.set('titleScratch', model.get('title'));
    controller.set('bodyScratch', model.get('body'));
  },

  renderTemplate: function() {
    this.render('post/edit', { controller: 'post.edit' });
  }
});

export default PostsNewRoute;

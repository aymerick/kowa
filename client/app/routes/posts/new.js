import AuthenticatedRoute from 'kowa/routes/authenticated';
import ContentEditionRoute from 'kowa/mixins/content-edition-route';
import Post from 'kowa/models/post';

// Uses post/edit controller and template, but with a new record
var PostsNewRoute = AuthenticatedRoute.extend(ContentEditionRoute, {
  controllerName: 'post.edit',

  model: function() {
    return this.store.createRecord('post', Post.newRecordAttrs({
      site: this.modelFor('site')
    }));
  },

  renderTemplate: function() {
    this.render('post/edit', { controller: 'post.edit' });
  }
});

export default PostsNewRoute;

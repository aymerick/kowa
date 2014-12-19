import AuthenticatedRoute from 'kowa/routes/authenticated';
import PostEditionRoute from 'kowa/mixins/post-edition-route';
import Post from 'kowa/models/post';

// Uses post/edit controller and template, but with a new record
var PostsNewRoute = AuthenticatedRoute.extend(PostEditionRoute, {
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

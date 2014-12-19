import AuthenticatedRoute from 'kowa/routes/authenticated';
import ContentEditionRoute from 'kowa/mixins/content-edition-route';
import Post from 'kowa/models/post';

var PostsNewRoute = AuthenticatedRoute.extend(ContentEditionRoute, {
  // use PostEditController
  controllerName: 'post.edit',

  // this is a fresh new model
  model: function() {
    return this.store.createRecord('post', Post.newRecordAttrs({
      site: this.modelFor('site')
    }));
  },

  // use 'post/edit' template
  renderTemplate: function() {
    this.render('post/edit');
  }
});

export default PostsNewRoute;

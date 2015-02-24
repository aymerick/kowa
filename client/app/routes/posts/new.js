import AuthenticatedRoute from 'kowa/routes/authenticated';
import ContentEditionRoute from 'kowa/mixins/content-edition-route';

var PostsNewRoute = AuthenticatedRoute.extend(ContentEditionRoute, {
  // use PostsPostController
  controllerName: 'posts.post',

  // this is a fresh new model
  model: function() {
    return this.store.createRecord('post', { site: this.modelFor('site') });
  },

  // use 'posts/post' template
  renderTemplate: function() {
    this.render('posts/post');
  }
});

export default PostsNewRoute;

import Ember from 'ember';

var PostController = Ember.ObjectController.extend({
    actions: {
        showPostContent: function () {
            this.transitionToRoute('posts.post', this.get('model'));
        }
    }
});

export default PostController;

import Ember from 'ember';
import PaginationControllerMixin from 'kowa/mixins/pagination-controller';

var PostsController = Ember.ArrayController.extend(PaginationControllerMixin, {
  sortProperties: ['createdAt'],
  sortAscending: false,

  // i18n for attributes values and components parameters
  i18n: function() {
    return {
      newPost: this.t('post.new'),
      editPost: this.t('post.edit'),
      coverImage: this.t('coverImage')
    };
  }.property() // @todo Watch current language
});

export default PostsController;

import Ember from 'ember';
import PaginationControllerMixin from 'kowa/mixins/pagination-controller';

var PostsController = Ember.ArrayController.extend(PaginationControllerMixin, {
    sortProperties: ['createdAt'],
});

export default PostsController;

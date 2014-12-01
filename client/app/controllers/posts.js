import Ember from 'ember';
import PaginationControllerMixin from 'kowa/mixins/pagination-controller';

var PostsController = Ember.ArrayController.extend(PaginationControllerMixin, {

});

export default PostsController;

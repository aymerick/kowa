import Ember from 'ember';
import PaginationControllerMixin from 'kowa/mixins/pagination-controller';

var ImagesController = Ember.ArrayController.extend(PaginationControllerMixin, {
    sortProperties: ['createdAt'],
    sortAscending: false
});

export default ImagesController;

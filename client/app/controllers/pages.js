import Ember from 'ember';
import PaginationControllerMixin from 'kowa/mixins/pagination-controller';

var PagesController = Ember.ArrayController.extend(PaginationControllerMixin, {
    sortProperties: ['createdAt'],
    sortAscending: false
});

export default PagesController;

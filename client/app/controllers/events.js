import Ember from 'ember';
import PaginationControllerMixin from 'kowa/mixins/pagination-controller';

var EventsController = Ember.ArrayController.extend(PaginationControllerMixin, {
  sortProperties: ['startAt'],
  sortAscending: false
});

export default EventsController;

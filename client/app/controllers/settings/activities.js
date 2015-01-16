import Ember from 'ember';
import PaginationControllerMixin from 'kowa/mixins/pagination-controller';

var SettingsActivitiesController = Ember.ArrayController.extend(PaginationControllerMixin, {
    sortProperties: ['createdAt'],
    sortAscending: false
});

export default SettingsActivitiesController;

import Ember from 'ember';
import PaginationControllerMixin from 'kowa/mixins/pagination-controller';

var ImagesController = Ember.ArrayController.extend(PaginationControllerMixin, {
    sortProperties: ['createdAt'],
    sortAscending: false,
    itemController: 'settings/image',

    needs: ['settings'],
    site: Ember.computed.alias('controllers.settings.model'),

    uploadUrl: function() {
      return "/api/images/upload?site=" + this.get('site').get('id')
    }.property('site')
});

export default ImagesController;

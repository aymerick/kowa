import Ember from 'ember';
import PaginationControllerMixin from 'kowa/mixins/pagination-controller';

var PagesController = Ember.ArrayController.extend(PaginationControllerMixin, {
  sortProperties: ['createdAt'],

  // i18n for attributes values and components parameters
  i18n: function() {
    return {
      newPage: this.t('page.new'),
      editPage: this.t('page.edit'),
      coverImage: this.t('coverImage')
    };
  }.property() // @todo Watch current language
});

export default PagesController;

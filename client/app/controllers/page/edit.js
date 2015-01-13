import Ember from 'ember';
import Page from 'kowa/models/page';

// @todo Refactor that into a ContentEditionController Mixin
var PageEditController = Ember.ObjectController.extend({
  isDirty: function() {
    var model = this.get('model');
    return ((model.get('isNew')) ||
            (model.get('title') !== this.get('titleScratch')) ||
            (model.get('tagline') !== this.get('taglineScratch')) ||
            (model.get('body')  !== this.get('bodyScratch')) ||
            (model.get('cover') !== this.get('coverScratch')));
  }.property('titleScratch', 'taglineScratch', 'bodyScratch', 'coverScratch', 'model.title', 'model.tagline', 'model.body', 'model.cover', 'model.isNew'),

  nothingChanged: Ember.computed.not('isDirty'),

  setupEdition: function() {
    var model = this.get('model');

    this.set('titleScratch', model.get('title'));
    this.set('taglineScratch', model.get('tagline'));
    this.set('bodyScratch', model.get('body'));
    this.set('coverScratch', model.get('cover'));
  },

  actions: {
    // called by 'select-image' modal controller
    imageSelected: function(field, image) {
      this.set(field, image);
    },

    removeCover: function() {
      this.set('coverScratch', null);
    },

    savePage: function() {
      if (!this.get('isDirty')) {
        // This should never happen
        return;
      }

      // set a default title
      if (!this.get('titleScratch')) {
          this.set('titleScratch', '(Untitled)');
      }

      // update model
      var model = this.get('model');
      model.set('title', this.get('titleScratch'));
      model.set('tagline', this.get('taglineScratch'));
      model.set('body', this.get('bodyScratch'));
      model.set('cover', this.get('coverScratch'));

      // persist on server
      var self = this;
      model.save().then(function (pageSaved) {
        self.get('flashes').success('Page saved.');

        return pageSaved;
      }).catch(function () {
        self.get('flashes').danger('Failed to save page.');

        // rollback model
        if (model.get('isNew')) {
          model.setProperties(Page.newRecordAttrs());
        }
        else {
          model.rollback();
        }
      });
    }
  }
});

export default PageEditController;

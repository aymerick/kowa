import Ember from 'ember';
import EmberUploader from 'ember-uploader/uploader';

import PaginationControllerMixin from 'kowa/mixins/pagination-controller';

var ImagesController = Ember.ArrayController.extend(PaginationControllerMixin, {
  sortProperties: ['createdAt'],
  sortAscending: false,
  itemController: 'settings/image',

  needs: ['settings'],
  site: Ember.computed.alias('controllers.settings.model'),

  uploadUrl: function() {
    return "/api/images/upload?site=" + this.get('site').get('id');
  }.property('site'),

  actions: {
    uploadImage: function(file) {
      var uploader = EmberUploader.create({ url: this.get('uploadUrl') });

      var now = new Date();
      var self = this;

      // create local record
      var image = this.store.createRecord('image', {
        createdAt: now,
        updatedAt: now,
        name: file.name,
        size: file.size,
        type: file.type,
        isUploading: true,
        uploadProgress: 0
      });

      uploader.on('progress', function(e) {
        image.set('uploadProgress', e.percent);
      });

      uploader.on('didUpload', function(response) {
        // update record
        var recordData = self.store.normalize('image', response.image);
        image.setupData(recordData);
        image.materializeId(recordData.id);

        // upload is over
        image.set('isUploading', false);
        image.set('uploadProgress', null);

        // remove record dirtyness
        image.transitionTo('loaded.saved');
      });

      // @todo on error !

      uploader.upload(file);
    }
  }
});

export default ImagesController;

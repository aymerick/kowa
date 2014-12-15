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

      // create temporary image
      var tmpImage = this.store.createRecord('image', {
        createdAt: now,
        updatedAt: now,
        name: file.name,
        size: file.size,
        type: file.type,
        isUploading: true,
        uploadProgress: 0
      });

      uploader.on('progress', function(e) {
        tmpImage.set('uploadProgress', e.percent);
      });

      uploader.on('didUpload', function(response) {
        //
        // This is crazy but I don't find any correct way to just modify the localy created record and set its 'isDirty' and 'isNew' properties to false.
        // So we destroy it, and recreate a new one with the payload sent back by the server.
        //
        // cf. http://stackoverflow.com/questions/13342250/how-to-manually-set-an-object-state-to-clean-saved-using-ember-data/27482276
        //

        // destroy temporary image
        tmpImage.destroyRecord();

        // push uploaded image
        var recordData = self.store.normalize('image', response.image);
        self.store.update('image', recordData);
      });

      // @todo on error !

      uploader.upload(file);
    }
  }
});

export default ImagesController;

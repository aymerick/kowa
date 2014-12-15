import Ember from 'ember';

import PaginationControllerMixin from 'kowa/mixins/pagination-controller';
import UploadJob from 'kowa/services/upload-job';

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
      var job = new UploadJob({
        url: this.get('uploadUrl'),
        file: file
      });

      // set dummy date to preserve images order when upload is finished and correct date is sent by the server
      var at = null;
      var lastDate = this.get('firstObject').get('createdAt');
      if (lastDate) {
        at = moment(lastDate).add(1, 'hours').toDate();
      }
      else {
        at = new Date();
      }

      // create temporary image
      var tmpImage = this.store.createRecord('image', {
        createdAt: at,
        updatedAt: at,
        name: file.name,
        size: file.size,
        type: file.type,
        isUploading: true,
        uploadProgress: 0
      });

      var controller = this;

      job.on('willUpload', function() {
        job.set('tmpImage', tmpImage);
      });

      job.on('progress', function(val) {
        job.get('tmpImage').set('uploadProgress', val);
      });

      job.on('didUpload', function(response) {
        //
        // This is crazy but I don't find any correct way to just modify the localy created record and set its 'isDirty' and 'isNew' properties to false.
        // So we destroy it, and recreate a new one with the payload sent back by the server.
        //
        // cf. http://stackoverflow.com/questions/13342250/how-to-manually-set-an-object-state-to-clean-saved-using-ember-data/27482276
        //

        // destroy temporary image
        job.get('tmpImage').destroyRecord();

        // push uploaded image
        var recordData = controller.store.normalize('image', response.image);
        controller.store.update('image', recordData);
      });

      this.uploader.schedule(job);
    }
  }
});

export default ImagesController;

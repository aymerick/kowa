import Ember from 'ember';
import moment from 'moment';

import PaginationControllerMixin from 'kowa/mixins/pagination-controller';
import UploadJob from 'kowa/services/upload-job';

var SettingsImagesController = Ember.ArrayController.extend(PaginationControllerMixin, {
  sortProperties: ['createdAt'],
  sortAscending: false,
  itemController: 'settings/image',

  needs: ['settings'],
  site: Ember.computed.alias('controllers.settings.model'),

  uploadUrl: function() {
    return "/api/images/upload?site=" + this.get('site').get('id');
  }.property('site'),

  // compute a dummy date to preserve images order when upload is finished and correct date is sent by the server
  // index is the file index when uploading multiple files at once
  computeTmpImageDate: function(index) {
    var firstImage = this.get('firstObject');
    if (firstImage) {
      var lastDate = firstImage.get('createdAt');
      if (lastDate) {
        return moment(lastDate).add(index + 1, 'hours').toDate();
      } else {
        return new Date();
      }
    } else {
      return new Date();
    }
  },

  actions: {
    uploadImage: function(file, index) {
      var job = new UploadJob({
        url: this.get('uploadUrl'),
        file: file
      });

      var at = this.computeTmpImageDate(index);

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

      job.on('uploadProgress', function(val) {
        job.get('tmpImage').set('uploadProgress', val);
      });

      job.on('didUpload', function(response) {
        //
        // This is crazy but I don't find any correct way to just modify the localy created record and set its 'isDirty' and 'isNew' properties to false.
        // So we destroy it, and recreate a new one with the payload sent back by the server.
        //
        // cf. http://stackoverflow.com/questions/13342250/how-to-manually-set-an-object-state-to-clean-saved-using-ember-data/27482276
        //
        job.get('tmpImage').destroyRecord();

        // push uploaded image
        var recordData = controller.store.normalize('image', response.image);
        controller.store.push('image', recordData);
      });

      this.uploader.schedule(job);
    }
  }
});

export default SettingsImagesController;

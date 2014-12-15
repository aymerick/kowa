import UploaderService from 'kowa/services/uploader';

var UploaderInitializer = {
    name: 'injectUploaderService',

    initialize: function (container, application) {
        application.register('uploader:main', UploaderService);

        application.inject('controller', 'uploader', 'uploader:main');
    }
};

export default UploaderInitializer;

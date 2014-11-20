// source: https://github.com/TryGhost/Ghost/blob/master/core/client/routes/error404.js

import Ember from 'ember';

var Error404Route = Ember.Route.extend({
    templateName: 'error',

    model: function () {
        return {
            status: 404
        };
    }
});

export default Error404Route;

import Ember from 'ember';

var PostEditionRoute = Ember.Mixin.create({
    setupController: function (controller, model) {
        this._super(controller, model);

        // init temp data
        controller.set('titleScratch', model.get('title'));
        controller.set('bodyScratch', model.get('body'));
    }
});

export default PostEditionRoute;

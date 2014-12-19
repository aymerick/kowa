import Ember from 'ember';

var PostEditionRoute = Ember.Mixin.create({
  setupController: function (controller, model) {
    this._super(controller, model);

    // init edition
    controller.resetEdition();
  }
});

export default PostEditionRoute;

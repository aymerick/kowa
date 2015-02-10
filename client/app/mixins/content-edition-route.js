import Ember from 'ember';

var ContentEditionRoute = Ember.Mixin.create({
  setupController: function (controller, model) {
    this._super(controller, model);

    // setup edition
    controller.setupEdition(controller.get('editionRelationships'));
  },

  actions: {
    willTransition: function (transition) {
      var controller = this.get('controller');
      if (controller.get('isDirty')) {
        transition.abort();
        this.send('openModal', 'leave-edition', controller, transition);
        return;
      }
    }
  }
});

export default ContentEditionRoute;

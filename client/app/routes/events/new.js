import AuthenticatedRoute from 'kowa/routes/authenticated';
import ContentEditionRoute from 'kowa/mixins/content-edition-route';

var EventsNewRoute = AuthenticatedRoute.extend(ContentEditionRoute, {
  // use EventEditController
  controllerName: 'event.edit',

  // this is a fresh new model
  model: function() {
    return this.store.createRecord('event', { site: this.modelFor('site') });
  },

  // use 'event/edit' template
  renderTemplate: function() {
    this.render('event/edit');
  }
});

export default EventsNewRoute;

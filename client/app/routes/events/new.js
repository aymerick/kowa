import AuthenticatedRoute from 'kowa/routes/authenticated';
import ContentEditionRoute from 'kowa/mixins/content-edition-route';

var EventsNewRoute = AuthenticatedRoute.extend(ContentEditionRoute, {
  // use EventsEventController
  controllerName: 'events.event',

  // this is a fresh new model
  model: function() {
    return this.store.createRecord('event', { site: this.modelFor('site') });
  },

  // use event' template
  renderTemplate: function() {
    this.render('events/event');
  }
});

export default EventsNewRoute;

import AuthenticatedRoute from 'kowa/routes/authenticated';
import ContentEditionRoute from 'kowa/mixins/content-edition-route';
import EventModel from 'kowa/models/event';

var EventsNewRoute = AuthenticatedRoute.extend(ContentEditionRoute, {
  // use EventEditController
  controllerName: 'event.edit',

  // this is a fresh new model
  model: function() {
    return this.store.createRecord('event', EventModel.newRecordAttrs({
      site: this.modelFor('site')
    }));
  },

  // use 'event/edit' template
  renderTemplate: function() {
    this.render('event/edit');
  }
});

export default EventsNewRoute;

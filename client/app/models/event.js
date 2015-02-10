import DS from 'ember-data';
import Ember from 'ember';

var EventModel = DS.Model.extend({
  createdAt: DS.attr('date'),
  updatedAt: DS.attr('date'),

  startDate: DS.attr('date'),
  endDate: DS.attr('date'),
  title: DS.attr(),
  body: DS.attr(),
  place: DS.attr(),

  cover: DS.belongsTo('image', { inverse: 'coverEvents', async: true }),

  site: DS.belongsTo('site', { async: true })
});

EventModel.reopenClass({
  newRecordAttrs: function(moreAttrs) {
    var now = new Date();

    return Ember.merge({
      createdAt: now,
      updatedAt: now,
      title: "",
      body: "",
      place: ""
    }, moreAttrs || { });
  }
});

export default EventModel;

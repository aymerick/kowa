/* global moment */

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

  site: DS.belongsTo('site', { async: true }),

  dateSentence: function() {
    return moment(this.get('startDate')).twix(moment(this.get('endDate'))).format();
  }.property('startDate', 'endDate')
});

EventModel.reopenClass({
  newRecordAttrs: function(moreAttrs) {
    var now = new Date();
    var endDate = moment(now).add(1, 'hours').toDate();

    return Ember.merge({
      createdAt: now,
      updatedAt: now,
      startDate: now,
      endDate: endDate,
      title: "",
      body: "",
      place: ""
    }, moreAttrs || { });
  }
});

export default EventModel;

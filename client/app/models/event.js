/* global moment */
import DS from 'ember-data';
import Ember from 'ember';

var EventModel = DS.Model.extend({
  createdAt: DS.attr('moment-date', { defaultValue: function() { return moment(new Date()); } }),
  updatedAt: DS.attr('moment-date', { defaultValue: function() { return moment(new Date()); } }),

  startDate: DS.attr('moment-date', { defaultValue: function() { return moment(new Date()).add(1, 'hours'); } }),
  endDate: DS.attr('moment-date', { defaultValue: function() { return moment(new Date()).add(2, 'hours'); } }),

  title: DS.attr('string', { defaultValue: "" }),
  body: DS.attr('string', { defaultValue: "" }),
  format: DS.attr('string', { defaultValue: "html" }), // Keep in sync with models.DEFAULT_FORMAT
  place: DS.attr('string', { defaultValue: "" }),

  cover: DS.belongsTo('image', { inverse: 'coverEvents', async: true }),

  site: DS.belongsTo('site', { async: true }),

  isMarkdown: Ember.computed.equal('format', 'md'), // Keep in sync with models.FORMAT_MARKDOWN

  dateSentence: function() {
    return moment(this.get('startDate')).twix(moment(this.get('endDate'))).format();
  }.property('startDate', 'endDate')
});

export default EventModel;

/* global moment */
import DS from 'ember-data';
import Ember from 'ember';

var Page = DS.Model.extend({
  createdAt: DS.attr('moment-date', { defaultValue: function() { return moment(new Date()); } }),
  updatedAt: DS.attr('moment-date', { defaultValue: function() { return moment(new Date()); } }),

  title: DS.attr('string', { defaultValue: "" }),
  tagline: DS.attr('string', { defaultValue: "" }),
  body: DS.attr('string', { defaultValue: "" }),
  format: DS.attr('string', { defaultValue: "html" }), // Keep in sync with models.DEFAULT_FORMAT

  inNavBar: DS.attr('boolean', { defaultValue: false }),

  cover: DS.belongsTo('image', { inverse: 'coverPages', async: true }),

  site: DS.belongsTo('site', { async: true }),

  isMarkdown: Ember.computed.equal('format', 'md') // Keep in sync with models.FORMAT_MARKDOWN
});

export default Page;

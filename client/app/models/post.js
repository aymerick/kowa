/* global moment */
import DS from 'ember-data';
import Ember from 'ember';

var Post = DS.Model.extend({
  createdAt: DS.attr('moment-date', { defaultValue: function() { return moment(new Date()); } }),
  updatedAt: DS.attr('moment-date', { defaultValue: function() { return moment(new Date()); } }),

  publishedAt: DS.attr('moment-date'),
  title: DS.attr('string', { defaultValue: "" }),
  body: DS.attr('string', { defaultValue: "" }),
  format: DS.attr('string', { defaultValue: "html" }), // Keep in sync with models.DEFAULT_FORMAT

  cover: DS.belongsTo('image', { inverse: 'coverPosts', async: true }),

  site: DS.belongsTo('site', { async: true }),

  isMarkdown: Ember.computed.equal('format', 'md') // Keep in sync with models.FORMAT_MARKDOWN
});

export default Post;

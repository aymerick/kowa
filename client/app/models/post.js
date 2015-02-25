import DS from 'ember-data';

var Post = DS.Model.extend({
  createdAt: DS.attr('date', { defaultValue: function() { return new Date(); } }),
  updatedAt: DS.attr('date', { defaultValue: function() { return new Date(); } }),

  publishedAt: DS.attr('date'),
  title: DS.attr('string', { defaultValue: "" }),
  body: DS.attr('string', { defaultValue: "" }),
  format: DS.attr('string', { defaultValue: "html" }), // Keep in sync with models.DEFAULT_FORMAT

  cover: DS.belongsTo('image', { inverse: 'coverPosts', async: true }),

  site: DS.belongsTo('site', { async: true }),

  isMarkdown: Ember.computed.equal('format', 'md') // Keep in sync with models.FORMAT_MARKDOWN
});

export default Post;

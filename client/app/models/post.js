import DS from 'ember-data';

var Post = DS.Model.extend({
  createdAt: DS.attr('date', { defaultValue: function() { return new Date(); } }),
  updatedAt: DS.attr('date', { defaultValue: function() { return new Date(); } }),

  publishedAt: DS.attr('date'),
  title: DS.attr('string', { defaultValue: "" }),
  body: DS.attr('string', { defaultValue: "" }),

  cover: DS.belongsTo('image', { inverse: 'coverPosts', async: true }),

  site: DS.belongsTo('site', { async: true })
});

export default Post;

import DS from 'ember-data';

var Post = DS.Model.extend({
  createdAt: DS.attr('date'),
  updatedAt: DS.attr('date'),

  publishedAt: DS.attr('date'),
  title: DS.attr(''),
  body: DS.attr(''),

  site: DS.belongsTo('site', { async: true })
});

export default Post;

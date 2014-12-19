import DS from 'ember-data';

var Post = DS.Model.extend({
  createdAt: DS.attr('date'),
  updatedAt: DS.attr('date'),

  publishedAt: DS.attr('date'),
  title: DS.attr(''),
  body: DS.attr(''),

  site: DS.belongsTo('site', { async: true })
});

Post.reopenClass({
  newRecordAttrs: function() {
    var now = new Date();

    return {
      createdAt: now,
      updatedAt: now,
      title: "(Untitled)",
      body: ""
    }
  }
});

export default Post;

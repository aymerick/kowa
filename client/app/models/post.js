import DS from 'ember-data';
import Ember from 'ember';

var Post = DS.Model.extend({
  createdAt: DS.attr('date'),
  updatedAt: DS.attr('date'),

  publishedAt: DS.attr('date'),
  title: DS.attr(''),
  body: DS.attr(''),

  cover: DS.belongsTo('image', { inverse: 'coverPosts' }),

  site: DS.belongsTo('site', { async: true })
});

Post.reopenClass({
  newRecordAttrs: function(moreAttrs) {
    var now = new Date();

    return Ember.merge({
      createdAt: now,
      updatedAt: now,
      title: "(Untitled)",
      body: ""
    }, moreAttrs || { });
  }
});

export default Post;

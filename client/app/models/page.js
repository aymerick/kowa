import DS from 'ember-data';
import Ember from 'ember';

var Page = DS.Model.extend({
  createdAt: DS.attr('date'),
  updatedAt: DS.attr('date'),

  title: DS.attr(''),
  tagline: DS.attr(''),
  body: DS.attr(''),

  cover: DS.belongsTo('image', { inverse: 'coverPages' }),

  site: DS.belongsTo('site', { async: true })
});

Page.reopenClass({
  newRecordAttrs: function(moreAttrs) {
    var now = new Date();

    return Ember.merge({
      createdAt: now,
      updatedAt: now,
      title: "",
      tagline: "",
      body: ""
    }, moreAttrs || { });
  }
});

export default Page;

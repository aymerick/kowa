import DS from 'ember-data';
import Ember from 'ember';

var Activity = DS.Model.extend({
  createdAt: DS.attr('date'),
  updatedAt: DS.attr('date'),

  title: DS.attr(),
  summary: DS.attr(),
  body: DS.attr(),

  cover: DS.belongsTo('image', { inverse: 'coverActivities', async: true }),

  site: DS.belongsTo('site', { async: true })
});

Activity.reopenClass({
  newRecordAttrs: function(moreAttrs) {
    var now = new Date();

    return Ember.merge({
      createdAt: now,
      updatedAt: now,
      title: "",
      summary: "",
      body: ""
    }, moreAttrs || { });
  }
});

export default Activity;

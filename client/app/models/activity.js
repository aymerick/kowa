import DS from 'ember-data';

var Activity = DS.Model.extend({
  createdAt: DS.attr('date', { defaultValue: function() { return new Date(); } }),
  updatedAt: DS.attr('date', { defaultValue: function() { return new Date(); } }),

  title: DS.attr('string', { defaultValue: "" }),
  summary: DS.attr('string', { defaultValue: "" }),
  body: DS.attr('string', { defaultValue: "" }),

  cover: DS.belongsTo('image', { inverse: 'coverActivities', async: true }),

  site: DS.belongsTo('site', { async: true })
});

export default Activity;

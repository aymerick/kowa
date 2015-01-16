import DS from 'ember-data';

var Activity = DS.Model.extend({
  createdAt: DS.attr('date'),
  updatedAt: DS.attr('date'),

  title: DS.attr(''),
  body: DS.attr(''),

  cover: DS.belongsTo('image', { inverse: 'coverActivities' }),

  site: DS.belongsTo('site', { async: true })
});

export default Activity;

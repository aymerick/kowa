import DS from 'ember-data';

var Site = DS.Model.extend({
  createdAt: DS.attr('date'),

  name: DS.attr(),
  tagline: DS.attr(),
  description: DS.attr(),
  moreDesc: DS.attr(),
  joinText: DS.attr(),

  user: DS.belongsTo('user', { async: true }),

  posts: DS.hasMany('post', { async: true }),
  // events: DS.hasMany('event', { async: true }),
  // pages: DS.hasMany('page', { async: true }),
  // actions: DS.hasMany('action', { async: true }),

  pageSettings: DS.hasMany('sitePageSetting')
});

export default Site;

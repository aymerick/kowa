import DS from 'ember-data';

var Site = DS.Model.extend({
  createdAt: DS.attr('date'),

  name: DS.attr(),
  tagline: DS.attr(),
  description: DS.attr(),
  moreDesc: DS.attr(),
  joinText: DS.attr(),

  logo: DS.belongsTo('image', { inverse: 'logoSites' }),
  cover: DS.belongsTo('image', { inverse: 'coverSites' }),

  user: DS.belongsTo('user', { async: true }),

  posts: DS.hasMany('post', { async: true }),
  // events: DS.hasMany('event', { async: true }),
  // pages: DS.hasMany('page', { async: true }),
  // activities: DS.hasMany('activity', { async: true }),
  images: DS.hasMany('image', { inverse: 'site', async: true }),

  pageSettings: DS.hasMany('sitePageSetting')
});

export default Site;

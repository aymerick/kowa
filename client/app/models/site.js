/* global moment */
import DS from 'ember-data';

var Site = DS.Model.extend({
  createdAt: DS.attr('moment-date'),
  updatedAt: DS.attr('moment-date'),
  changedAt: DS.attr('moment-date'),
  builtAt: DS.attr('moment-date'),

  name: DS.attr(),
  tagline: DS.attr(),
  description: DS.attr(),
  moreDesc: DS.attr(),
  joinText: DS.attr(),

  email: DS.attr(),
  address: DS.attr(),

  facebook: DS.attr(),
  twitter: DS.attr(),
  googlePlus: DS.attr(),

  logo: DS.belongsTo('image', { inverse: 'logoSites', async: true }),
  cover: DS.belongsTo('image', { inverse: 'coverSites', async: true }),

  user: DS.belongsTo('user', { async: true }),

  posts: DS.hasMany('post', { async: true }),
  events: DS.hasMany('event', { async: true }),
  pages: DS.hasMany('page', { async: true }),
  activities: DS.hasMany('activity', { async: true }),
  images: DS.hasMany('image', { inverse: 'site', async: true }),

  pageSettings: DS.hasMany('site-page-setting'),

  // settings
  theme: DS.attr(),
  uglyUrl: DS.attr('boolean')
});

export default Site;

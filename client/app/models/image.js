import DS from 'ember-data';

var ImageModel = DS.Model.extend({
  createdAt: DS.attr('date'),
  updatedAt: DS.attr('date'),
  name: DS.attr(),
  size: DS.attr('number'),
  type: DS.attr(),

  url: DS.attr(),
  smallUrl: DS.attr(),
  thumbUrl: DS.attr(),
  mediumUrl: DS.attr(),
  mediumCropUrl: DS.attr(),

  site: DS.belongsTo('site', { inverse: 'images', async: true }),

  logoSites: DS.hasMany('site', { inverse: 'logo', async: true }),
  coverSites: DS.hasMany('site', { inverse: 'cover', async: true }),
  coverPosts: DS.hasMany('post', { inverse: 'cover', async: true }),
  coverPages: DS.hasMany('page', { inverse: 'cover', async: true }),
  coverActivities: DS.hasMany('activity', { inverse: 'cover', async: true })
});

export default ImageModel;

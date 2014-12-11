import DS from 'ember-data';

var SITE_SERIALIZER_FIELDS = [ 'name', 'tagline', 'description', 'moreDesc', 'joinText', 'logo', 'cover' ];

var SiteSerializer = DS.RESTSerializer.extend({
  serialize: function(site, options) {
    var result = site.getProperties(SITE_SERIALIZER_FIELDS);

    if (options && options.includeId) {
      result['id'] = site.get('id');
    }

    // logo
    if (result['logo']) {
      result['logo'] = result['logo']['id'];
    }

    if (result['logo'] == null) {
      delete(result['logo']);
    }

    // cover
    if (result['cover']) {
      result['cover'] = result['cover']['id'];
    }

    if (result['cover'] == null) {
      delete(result['cover']);
    }

    return result;
  }
});

export default SiteSerializer;

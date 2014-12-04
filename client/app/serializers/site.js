import DS from 'ember-data';

var SITE_SERIALIZER_FIELDS = [ 'name', 'tagline', 'description', 'moreDesc', 'joinText' ];

var SiteSerializer = DS.RESTSerializer.extend({
  serialize: function(site, options) {
    var result = site.getProperties(SITE_SERIALIZER_FIELDS);

    if (options && options.includeId) {
      result['id'] = site.get('id');
    }

    return result;
  }
});

export default SiteSerializer;

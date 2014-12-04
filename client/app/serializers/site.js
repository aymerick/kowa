import DS from 'ember-data';

var SiteSerializer = DS.RESTSerializer.extend({
  serialize: function(site, options) {
    var fields = [ 'name', 'tagline', 'description', 'moreDesc', 'joinText' ];
    var result = { };

    fields.forEach(function(field) {
      result[field] = site.get(field);
    });

    if (options && options.includeId) {
      result['id'] = site.get('id');
    }

    return result;
  }
});

export default SiteSerializer;

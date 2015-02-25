import DS from 'ember-data';

var PAGE_SERIALIZER_FIELDS = [ 'title', 'tagline', 'body', 'format', 'site', 'cover', 'inNavBar' ];

var PageSerializer = DS.RESTSerializer.extend({
  serialize: function(page, options) {
    var result = page.getProperties(PAGE_SERIALIZER_FIELDS);

    if (options && options.includeId) {
      var pageId = page.get('id');
      if (pageId != null) {
        result['id'] = pageId;
      }
    }

    // site
    if (result['site']) {
      result['site'] = result['site'].get('id');
    }

    // cover
    if (result['cover']) {
      result['cover'] = result['cover'].get('id');
    }

    if (result['cover'] == null) {
      delete(result['cover']);
    }

    return result;
  }
});

export default PageSerializer;

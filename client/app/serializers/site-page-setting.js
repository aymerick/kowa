import DS from 'ember-data';

var SITE_PAGE_SETTINGS_SERIALIZER_FIELDS = [
  'kind', 'title', 'tagline', 'cover', 'disabled'
];

var SitePageSettingsSerializer = DS.RESTSerializer.extend({
  serialize: function(pageSettings, options) {
    var result = pageSettings.getProperties(SITE_PAGE_SETTINGS_SERIALIZER_FIELDS);

    if (options && options.includeId) {
      var settingsId = pageSettings.get('id');
      if (settingsId != null) {
        result['id'] = settingsId;
      }
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

export default SitePageSettingsSerializer;

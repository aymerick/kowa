import DS from 'ember-data';

var ACTIVITY_SERIALIZER_FIELDS = [ 'title', 'body', 'site', 'cover' ];

var ActivitySerializer = DS.RESTSerializer.extend({
  serialize: function(activity, options) {
    var result = activity.getProperties(ACTIVITY_SERIALIZER_FIELDS);

    if (options && options.includeId) {
      var activityId = activity.get('id');
      if (activityId != null) {
        result['id'] = activityId;
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

export default ActivitySerializer;

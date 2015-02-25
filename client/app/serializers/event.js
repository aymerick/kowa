import DS from 'ember-data';

var EVENT_SERIALIZER_FIELDS = [ 'startDate', 'endDate', 'title', 'body', 'format', 'place', 'site', 'cover' ];

var EventSerializer = DS.RESTSerializer.extend({
  serialize: function(eventModel, options) {
    var result = eventModel.getProperties(EVENT_SERIALIZER_FIELDS);

    if (options && options.includeId) {
      var eventId = eventModel.get('id');
      if (eventId != null) {
        result['id'] = eventId;
      }
    }

    // startDate
    if (result['startDate'] == null) {
      delete(result['startDate']);
    }

    // endDate
    if (result['endDate'] == null) {
      delete(result['endDate']);
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

export default EventSerializer;

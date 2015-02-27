import DS from 'ember-data';

var MEMBER_SERIALIZER_FIELDS = [ 'fullname', 'role', 'description', 'site', 'photo' ];

var MemberSerializer = DS.RESTSerializer.extend({
  serialize: function(member, options) {
    var result = member.getProperties(MEMBER_SERIALIZER_FIELDS);

    if (options && options.includeId) {
      var memberId = member.get('id');
      if (memberId != null) {
        result['id'] = memberId;
      }
    }

    // site
    if (result['site']) {
      result['site'] = result['site'].get('id');
    }

    // photo
    if (result['photo']) {
      result['photo'] = result['photo'].get('id');
    }

    if (result['photo'] == null) {
      delete(result['photo']);
    }

    return result;
  }
});

export default MemberSerializer;

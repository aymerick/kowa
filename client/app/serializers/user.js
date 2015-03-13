import DS from 'ember-data';

var USER_SERIALIZER_FIELDS = [
  'firstName', 'lastName', 'lang'
];

var UserSerializer = DS.RESTSerializer.extend({
  serialize: function(user, options) {
    var result = user.getProperties(USER_SERIALIZER_FIELDS);

    if (options && options.includeId) {
      result['id'] = user.get('id');
    }

    return result;
  }
});

export default UserSerializer;

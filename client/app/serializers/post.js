import DS from 'ember-data';

var POST_SERIALIZER_FIELDS = [ 'title', 'body' ];

var PostSerializer = DS.RESTSerializer.extend({
  serialize: function(post, options) {
    var result = post.getProperties(POST_SERIALIZER_FIELDS);

    if (options && options.includeId) {
      result['id'] = post.get('id');
    }

    return result;
  }
});

export default PostSerializer;

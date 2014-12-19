import DS from 'ember-data';

var POST_SERIALIZER_FIELDS = [ 'title', 'body', 'site' ];

var PostSerializer = DS.RESTSerializer.extend({
  serialize: function(post, options) {
    var result = post.getProperties(POST_SERIALIZER_FIELDS);

    if (options && options.includeId) {
      var postId = post.get('id');
      if (postId != null) {
        result['id'] = postId;
      }
    }

    // site
    if (result['site']) {
      result['site'] = result['site'].get('id');
    }

    return result;
  }
});

export default PostSerializer;

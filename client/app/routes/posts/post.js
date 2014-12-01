import Ember from 'ember';
import AuthenticatedRoute from 'kowa/routes/authenticated';

var PostRoute = AuthenticatedRoute.extend({
  model: function(params) {
    return this.store.find('post', params.post_id);
  }
});

export default PostRoute;

import AuthenticatedRoute from 'kowa/routes/authenticated';

var SitesRoute = AuthenticatedRoute.extend({
  model: function() {
    return this.session.get('currentUser').then(function(currentUser){
      return currentUser.get('sites');
    });
  }
});

export default SitesRoute;

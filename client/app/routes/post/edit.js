import AuthenticatedRoute from 'kowa/routes/authenticated';

var PostEditRoute = AuthenticatedRoute.extend({
  setupController: function (controller, model) {
    this._super(controller, model);

    // init temp data
    controller.set('titleScratch', model.get('title') || "Untitled");
    controller.set('bodyScratch', model.get('body'));
  }
});

export default PostEditRoute;

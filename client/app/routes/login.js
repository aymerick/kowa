import UnauthenticatedRoute from 'kowa/routes/unauthenticated';
import StyleBodyMixin from 'kowa/mixins/style-body';

var LoginRoute = UnauthenticatedRoute.extend(StyleBodyMixin, {
  classNames: ['kowa-login']
});

export default LoginRoute;

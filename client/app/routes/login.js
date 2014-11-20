import UnauthenticatedRoute from 'kowa/routes/unauthenticated';
import styleBody from 'kowa/mixins/style-body';

var LoginRoute = UnauthenticatedRoute.extend(styleBody, {
  classNames: ['kowa-login'],
});

export default LoginRoute;

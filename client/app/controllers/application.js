var ApplicationController = Ember.Controller.extend({
    hideNav: Ember.computed.match('currentPath', /(error|login|forgotten)/)
});

export default ApplicationController;

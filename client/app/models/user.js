import DS from 'ember-data';

var User = DS.Model.extend({
  createdAt: DS.attr('moment-date'),

  email: DS.attr(),
  firstName: DS.attr(),
  lastName: DS.attr(),

  sites: DS.hasMany('site', { async: true }),

  fullName: function() {
    return this.get('firstName') + ' ' + this.get('lastName');
  }.property('firstName', 'lastName')
});

export default User;

import DS from 'ember-data';

var User = DS.Model.extend({
  createdAt: DS.attr('date'),

  firstName: DS.attr(),
  lastName: DS.attr(),

  sites: DS.hasMany('site'),

  fullName: function() {
    return this.get('firstName') + ' ' + this.get('lastName');
  }.property('firstName', 'lastName')
});

export default User;

import DS from 'ember-data';

var Member = DS.Model.extend({
  createdAt: DS.attr('date', { defaultValue: function() { return new Date(); } }),
  updatedAt: DS.attr('date', { defaultValue: function() { return new Date(); } }),

  fullname: DS.attr('string', { defaultValue: "" }),
  role: DS.attr('string', { defaultValue: "" }),
  description: DS.attr('string', { defaultValue: "" }),

  photo: DS.belongsTo('image', { inverse: 'photoMembers', async: true }),

  site: DS.belongsTo('site', { async: true })
});

export default Member;

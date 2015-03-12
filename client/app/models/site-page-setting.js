import DS from 'ember-data';

export default DS.Model.extend({
  kind: DS.attr('string'),
  title: DS.attr('string', { defaultValue: "" }),
  tagline: DS.attr('string', { defaultValue: "" }),
  disabled: DS.attr('boolean'),

  // this fake attribute does not exist on the server but is needed
  // to build custom URLs when creating/updating a record
  site: DS.belongsTo('site'),

  cover: DS.belongsTo('image', { inverse: 'sitePageSettings', async: true })
});

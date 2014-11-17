import DS from 'ember-data';

export default DS.Model.extend({
  kind: DS.attr(),
  title: DS.attr(),
  tagline: DS.attr()
});

import DS from 'ember-data';

export default DS.Model.extend({
  created_at: DS.attr('date'),
  user_id: DS.attr(),

  name: DS.attr(),
  tagline: DS.attr(),
  description: DS.attr(),
  more_desc: DS.attr(),
  join_text: DS.attr(),
  page_settings: DS.hasMany('site-page-settings')
});

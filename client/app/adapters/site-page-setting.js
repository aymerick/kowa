import DS from 'ember-data';

var SitePageSettingAdapter = DS.RESTAdapter.extend({
  pathForType: function() {
    return 'page-settings';
  },

  buildURL: function(type, id, record) {
    // I know, side-effect programming is bad... I feel bad for that
    this.set('namespace', "api/sites/" + record.get('site').get('id'));

    return this._super(type, id, record);
  }
});

export default SitePageSettingAdapter;

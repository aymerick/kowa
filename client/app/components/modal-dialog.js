import Ember from 'ember';

var ModalDialog = Ember.Component.extend({
  closeText: 'Close',
  okText: 'Ok',

  actions: {
    ok: function() {
      this.$('.modal').modal('hide');
      this.sendAction('ok');
    }
  },

  show: function() {
    this.$('.modal').modal().on('hidden.bs.modal', function() {
      this.sendAction('close');
    }.bind(this));
  }.on('didInsertElement')
});

export default ModalDialog;

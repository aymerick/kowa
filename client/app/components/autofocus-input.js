import Ember from "ember";

var AutofocusInput = Ember.TextField.extend({
  focus: true,

  setFocus: function () {
    if (this.focus) {
      this.$().focus();
    }
  }.on('didInsertElement')
});

export default AutofocusInput;

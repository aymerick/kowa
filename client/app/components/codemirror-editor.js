import Ember from "ember";
import CodeMirror from "codemirror";

var CodemirrorEditor = Ember.TextArea.extend({
  value: null,

  // initialize codemirror when inserted
  initCodemirror: function() {
    var codeMirror = CodeMirror.fromTextArea(this.get('element'), {
      mode:'gfm',
      lineWrapping: true
    });

    this.set('codeMirror', codeMirror);

    // setup callbacks
    this.setupValueDidChange();
    this.setupCodeMirrorValueDidChange();

    this.valueDidChange();
  }.on('didInsertElement'),

  // refresh codemirror when visible
  refresh: function() {
    this.get('codeMirror').refresh();
  }.on('becameVisible'),

  // setup 'valueDidChange' callback
  setupValueDidChange: function() {
    this.addObserver('value', this, 'valueDidChange');

    // remove observer when destroyed
    this.on('willDestroyElement', this, function() {
      this.removeObserver('value', this, 'valueDidChange');
    });
  },

  // setup 'codeMirrorValueDidChange' callback
  setupCodeMirrorValueDidChange: function() {
    var callback = Ember.run.bind(this, 'codeMirrorValueDidChange');

    this.get('codeMirror').on('change', callback);

    // remove callback when destroyed
    this.on('willDestroyElement', this, function() {
      this.get('codeMirror').off('change', callback);
    });
  },

  // callback when value changed
  valueDidChange: function() {
    var codeMirror = this.get('codeMirror');
    var value      = this.get('value');

    if (value !== codeMirror.getValue()) {
      codeMirror.setValue(value || '');
    }
  },

  // callback when codemirror value changed
  codeMirrorValueDidChange: function(codeMirror) {
    this.set('value', codeMirror.getValue());
  }
});

export default CodemirrorEditor;

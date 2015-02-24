import Ember from "ember";
import CodeMirror from "codemirror";

var CodemirrorEditor = Ember.TextArea.extend({
  value: null,

  // initialize codemirror when inserted
  initCodemirror: function() {
    var codeMirror = CodeMirror.fromTextArea(this.get('element'), {
      mode:'gfm',
      lineWrapping: true,
      dragDrop: false
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
    var self = this;
    this.on('willDestroyElement', this, function() {
      self.removeObserver('value', self, 'valueDidChange');
    });
  },

  // setup 'codeMirrorValueDidChange' callback
  setupCodeMirrorValueDidChange: function() {
    var callback = Ember.run.bind(this, 'codeMirrorValueDidChange');

    this.get('codeMirror').on('change', callback);

    // remove callback when destroyed
    var self = this;
    this.on('willDestroyElement', this, function() {
      self.get('codeMirror').off('change', callback);
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

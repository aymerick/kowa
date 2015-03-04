import Ember from "ember";
import CodeMirror from "codemirror";

var CodemirrorEditor = Ember.TextArea.extend({
  editor: null,
  value: null,

  // initialize editor when inserted
  initEditor: function() {
    var editor = CodeMirror.fromTextArea(this.get('element'), {
      mode:'gfm',
      lineWrapping: true,
      dragDrop: false
    });

    this.set('editor', editor);

    // setup callbacks
    this.setupValueDidChange();
    this.setupEditorValueDidChange();

    this.valueDidChange();
  }.on('didInsertElement'),

  // refresh editor when visible
  refresh: function() {
    this.get('editor').refresh();
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

  // setup 'editorValueDidChange' callback
  setupEditorValueDidChange: function() {
    var callback = Ember.run.bind(this, 'editorValueDidChange');

    this.get('editor').on('change', callback);

    // remove callback when destroyed
    var self = this;
    this.on('willDestroyElement', this, function() {
      self.get('editor').off('change', callback);
    });
  },

  // callback when value changed
  valueDidChange: function() {
    var editor = this.get('editor');
    var value  = this.get('value');

    if (value !== editor.getValue()) {
      editor.setValue(value || '');
    }
  },

  // callback when value changed in editor
  editorValueDidChange: function(editor) {
    this.set('value', editor.getValue());
  }
});

export default CodemirrorEditor;

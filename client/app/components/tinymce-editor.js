import Ember from "ember";

var TinyMCEEditor = Ember.Component.extend({
  classNames: ['tinymce-editor'],
  classNameBindings: ['fillHeight'],

  value: null,
  editor: null,

  // settings
  height: null,

  // resize handling
  $window: null,
  resizeHandler: null,

  // callback to setup editor
  setupEditor: function(editor) {
    this.set('editor', editor);

    var self = this;

    // bind change event
    editor.on('change', function() {
      self.editorValueDidChange(editor.getContent());
    });
  },

  // initialize editor when inserted
  initEditor: function() {
    var settings = {
      theme_url: '/tinymce/themes/modern/theme.min.js',
      skin_url: '/tinymce/skins/lightgray',
      external_plugins: {
        hr: '/tinymce/plugins/hr/plugin.min.js',
        link: '/tinymce/plugins/link/plugin.min.js'
      },
      menubar: false,
      statusbar : false,
      resize: false,
      toolbar: 'styleselect | bold italic underline | alignleft aligncenter alignright | bullist numlist | link | hr',
      language_url: "/tinymce-locales/fr_FR.js", // @todo i18n
      setup: Ember.run.bind(this, this.setupEditor)
    };

    if (this.get('fillHeight')) {
      settings = Ember.merge(settings, {
        height: '100%'
      });

      this.resizeHandler = Ember.run.bind(this, this.onResize);

      this.$window = Ember.$(window);
      this.$window.on("resize", this.resizeHandler);
      this.resizeHandler();
    }
    else if (this.get('height') !== null) {
      settings = Ember.merge(settings, {
        height: this.get('height')
      });
    }

    this.$('textarea').tinymce(settings);

    this.setupValueDidChange();
  }.on('didInsertElement'),

  // remove editor when destroyed
  removeEditor: function(){
    if (this.resizeHandler !== null) {
      this.$window.off("resize", this.resizeHandler);
    }

    this.get('editor').destroy();
  }.on('willDestroyElement'),

  // setup 'valueDidChange' callback
  setupValueDidChange: function() {
    this.addObserver('value', this, 'valueDidChange');

    // remove observer when destroyed
    var self = this;
    this.on('willDestroyElement', this, function() {
      self.removeObserver('value', self, 'valueDidChange');
    });
  },

  // callback when value changed in editor
  editorValueDidChange: function(editedValue) {
    this.set('value', editedValue);
  },

  // callback when value changed
  valueDidChange: function() {
    var value = this.get('value');
    var editor = this.get('editor');

    if (value !== editor.getContent()) {
      debugger;
      editor.setContent(value);
    }
  },

  // callback when window size changed
  onResize: function() {
    var toolBarHeight = this.$('.mce-toolbar-grp').outerHeight();

    this.$('.mce-edit-area').css('top', toolBarHeight);
  }
});

export default TinyMCEEditor;
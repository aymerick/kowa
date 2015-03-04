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

    // bind keyup event
    editor.on('keyup', function() {
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
        link: '/tinymce/plugins/link/plugin.min.js',
        paste: '/tinymce/plugins/paste/plugin.min.js'
      },
      paste_as_text: true,
      menubar: false,
      statusbar : false,
      resize: false,
      toolbar: 'undo redo | styleselect | bold italic underline | alignleft aligncenter alignright | bullist numlist | link | hr',
      language_url: "/tinymce/locales/fr_FR.js", // @todo i18n
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

    // NOTE: Textarea value is set here instead of in the template because I don't want
    //       to use a bound value that causes edition hiccups.
    this.$('textarea').val(this.get('value'));
    this.$('textarea').tinymce(settings);
  }.on('didInsertElement'),

  // remove editor when destroyed
  removeEditor: function(){
    if (this.resizeHandler !== null) {
      this.$window.off("resize", this.resizeHandler);
    }

    this.get('editor').destroy();
  }.on('willDestroyElement'),

  // callback when value changed in editor
  editorValueDidChange: function(newValue) {
    this.sendAction('onChange', newValue);
  },

  // callback when window size changed
  onResize: function() {
    var toolBarHeight = this.$('.mce-toolbar-grp').outerHeight();

    this.$('.mce-edit-area').css('top', toolBarHeight);
  }
});

export default TinyMCEEditor;

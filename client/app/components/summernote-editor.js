import Ember from "ember";

var SummernoteEditor = Ember.Component.extend({
  classNames: ['summernote-editor'],
  classNameBindings: ['fillHeight'],

  value: null,

  // settings
  height: null,
  fillHeight: false,
  placeholder: false,

  // resize handling
  $window: null,
  resizeHandler: null,

  // initialize editor when inserted
  initEditor: function() {
    var settings = {
      placeholder: this.get('placeholder'),
      toolbar: [
        ['style', ['style']],
        ['font', ['bold', 'italic', 'underline', 'clear']],
        ['para', ['ul', 'ol', 'paragraph']],
        ['insert', ['link', 'hr']]
      ],
      styleTags: ['p', 'h1', 'h2', 'h3', 'h4', 'h5', 'blockquote'],
      disableDragAndDrop: true,
      disableResizeEditor: true,
      lang: 'fr-FR', // @todo FIXME i18n
      onChange: Ember.run.bind(this, this.editorValueDidChange)
    };

    if (this.get('height') !== null) {
      settings = Ember.merge(settings, {
        height: this.get('height'),
        minHeight: this.get('height')
      });
    }

    this.$('textarea').summernote(settings);

    if (this.get('fillHeight')) {
      this.resizeHandler = Ember.run.bind(this, this.onResize);

      this.$window = Ember.$(window);
      this.$window.on("resize", this.resizeHandler);
      this.resizeHandler();
    }

    this.setupValueDidChange();
  }.on('didInsertElement'),

  // remove editor when destroyed
  removeEditor: function(){
    if (this.resizeHandler !== null) {
      this.$window.off("resize", this.resizeHandler);
    }

    this.$('textarea').destroy();
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

    if (value !== this.$('.note-editable').html()) {
      this.$('textarea').code(value);
    }
  },

  // callback when window size changed
  onResize: function() {
    var topBarHeight = this.$('.note-editor .note-toolbar').outerHeight();

    this.$('.note-editor .note-editable').css('top', topBarHeight);
  }
});

export default SummernoteEditor;

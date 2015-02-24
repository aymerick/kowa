import Ember from "ember";

var SummernoteEditor = Ember.Component.extend({
  classNames: ['summernote-editor'],
  classNameBindings: ['fillHeight'],

  height: null,
  fillHeight: false,
  $window: null,
  resizeHandler: null,
  placeholder: false,

  // initialize summernote when inserted
  initSummernote: function() {
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
      onChange: Ember.run.bind(this, this.onChange)
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
  }.on('didInsertElement'),

  // remove summernote when destroyed
  removeSummernote: function(){
    if (this.resizeHandler !== null) {
      this.$window.off("resize", this.resizeHandler);
    }

    this.$('textarea').destroy();
  }.on('willDestroyElement'),

  // callback when content changed
  onChange: function(contents) {
    this.set('content', contents);
  },

  // callback when window size changed
  onResize: function() {
    var topBarHeight = this.$('.note-editor .note-toolbar').outerHeight();

    this.$('.note-editor .note-editable').css('top', topBarHeight);
  }
});

export default SummernoteEditor;

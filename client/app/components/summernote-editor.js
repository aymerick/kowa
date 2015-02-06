import Ember from "ember";

var SummernoteEditor = Ember.TextArea.extend({
  height: 300,
  placeholder: false,

  // initialize summernote when inserted
  initSummernote: function() {
    this.$().summernote({
      height: this.get('height'),
      placeholder: this.get('placeholder'),
      toolbar: [
        ['style', ['style']],
        ['font', ['bold', 'italic', 'underline', 'clear']],
        ['para', ['ul', 'ol', 'paragraph']],
        ['insert', ['link', 'picture', 'hr']]
      ],
      styleTags: ['p', 'h1', 'h2', 'h3', 'h4', 'h5', 'blockquote'],
      disableDragAndDrop: true,
      onChange: Ember.run.bind(this, this.onChange)
    });
  }.on('didInsertElement'),

  // remove summernote when destroyed
  removeSummernote: function(){
    this.$().destroy();
  }.on('willDestroyElement'),

  // callback when content changed
  onChange: function(contents) {
    this.set('value', contents);
  }
});

export default SummernoteEditor;

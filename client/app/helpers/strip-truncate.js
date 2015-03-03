import Ember from "ember";

export default Ember.Handlebars.makeBoundHelper(function(value, options) {
  var opts = options.hash;

  var start = opts.start || 0;
  var len = opts.max;

  // replace new lines with spaces
  value = value.replace(/\n/g, " ")
  value = value.replace(/<br\s*[\/]?>/gi, " ")

  // strip html tags
  var tmp = document.createElement("DIV");
  tmp.innerHTML = value;
  var strippedValue = tmp.textContent || tmp.innerText || "";

  // truncate
  var out = strippedValue.substr(start, len);

  if (strippedValue.length > len) {
    out += '...';
  }

  return out;
});
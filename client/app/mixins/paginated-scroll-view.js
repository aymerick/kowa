import Ember from 'ember';

var PaginatedScrollViewMixin = Ember.Mixin.create({
  attachCheckScroll: function () {
    var elt = this.$();
    elt.on('scroll', Ember.run.bind(this, this.checkScroll));
  }.on('didInsertElement'),

  detachCheckScroll: function () {
    var elt = this.$();
    elt.off('scroll');
  }.on('willDestroyElement'),

  checkScroll: function (event) {
    var elt          = event.target;
    var controller   = this.get('controller');
    var isLoading    = controller.get('isLoading');
    var triggerPoint = 100;

    // If we haven't passed our threshold or we are already fetching content, exit
    if (isLoading || (elt.scrollTop + elt.clientHeight + triggerPoint <= elt.scrollHeight)) {
      return;
    }

    controller.send('loadNextPage');
  }
});

export default PaginatedScrollViewMixin;

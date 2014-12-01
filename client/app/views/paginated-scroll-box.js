import Ember from 'ember';

var PaginatedScrollBox = Ember.View.extend({
    // setScrollClassName: function (options) {
    //     var el = this.$().closest('.content-list');

    //     var offset = 10;
    //     var className = 'scrolling';

    //     if (this.scrollTop() > offset) {
    //         el.addClass(className);
    //     } else {
    //         el.removeClass(className);
    //     }
    // },

    // attachScrollClassHandler: function () {
    //     var el = this.$();

    //     el.on('scroll', Ember.run.bind(this, this.setScrollClassName));
    // }.on('didInsertElement'),

    // detachScrollClassHandler: function () {
    //     this.$().off('scroll');
    // }.on('willDestroyElement'),

    attachCheckScroll: function () {
        var el = this.$();

        el.on('scroll', Ember.run.bind(this, this.checkScroll));
    }.on('didInsertElement'),

    detachCheckScroll: function () {
        var el = this.$();
        el.off('scroll');
    }.on('willDestroyElement'),

    checkScroll: function (event) {
        var element = event.target,
            triggerPoint = 100,
            controller = this.get('controller'),
            isLoading = controller.get('isLoading');

        // If we haven't passed our threshold or we are already fetching content, exit
        if (isLoading || (element.scrollTop + element.clientHeight + triggerPoint <= element.scrollHeight)) {
            return;
        }

        controller.send('loadNextPage');
    }
});

export default PaginatedScrollBox;

import Ember from 'ember';
import ItemView from 'kowa/views/item-view';

var PostItemView = ItemView.extend({
    // doubleClick: function () {
    //     this.get('controller').send('openEditor');
    // },

    click: function () {
        this.get('controller').send('showPostContent');
    }
});

export default PostItemView;

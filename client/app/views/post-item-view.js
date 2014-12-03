import ActiveItemView from 'kowa/views/active-item';

var PostItemView = ActiveItemView.extend({
    click: function () {
        this.get('controller').send('showPostContent');
    }
});

export default PostItemView;

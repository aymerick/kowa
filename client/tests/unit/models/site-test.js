import {
  moduleForModel,
  test
} from 'ember-qunit';

moduleForModel('site', 'Site', {
  // Specify the other units that are required for this test.
  needs: ['model:image', 'model:user', 'model:post', 'model:page', 'model:activity', 'model:site-page-setting', 'model:event']
});

test('it exists', function(assert) {
  var model = this.subject();
  // var store = this.store();
  assert.ok(!!model);
});

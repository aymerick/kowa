export default function(){
  // settings.contact
  this.transition(
    this.fromRoute('settings.contact.index'),
    this.toRoute('settings.contact.page'),
    this.use('toLeft'),
    this.reverse('toRight')
  );

  // settings.activities
  this.transition(
    this.fromRoute('settings.activities.index'),
    this.toRoute('settings.activities.activity'),
    this.use('toLeft'),
    this.reverse('toRight')
  );

  this.transition(
    this.fromRoute('settings.activities.index'),
    this.toRoute('settings.activities.page'),
    this.use('toLeft'),
    this.reverse('toRight')
  );

  // this.transition(
  //   this.fromRoute('settings.activities.index'),
  //   this.toRoute('settings.activities.activity'),
  //   this.use('explode', {
  //     pick: '.activity-cover',
  //     use: ['flyTo', { duration: 500 } ]
  //   }, { use: 'toLeft' }),
  //   this.reverse('explode', {
  //     pick: '.activity-cover',
  //     use: ['flyTo', { duration: 500 } ]
  //   }, { use: 'toRight' })
  // );

  // settings.members
  this.transition(
    this.fromRoute('settings.members.index'),
    this.toRoute('settings.members.member'),
    this.use('toLeft'),
    this.reverse('toRight')
  );

  this.transition(
    this.fromRoute('settings.members.index'),
    this.toRoute('settings.members.page'),
    this.use('toLeft'),
    this.reverse('toRight')
  );
}

export default function(){
  // settings.activities
  this.transition(
    this.fromRoute('settings.activities.index'),
    this.toRoute('settings.activities.activity'),
    this.use('toLeft'),
    this.reverse('toRight')
  );

  // settings.members
  this.transition(
    this.fromRoute('settings.members.index'),
    this.toRoute('settings.members.member'),
    this.use('toLeft'),
    this.reverse('toRight')
  );
}

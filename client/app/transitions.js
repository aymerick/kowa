export default function(){
  this.transition(
    this.fromRoute('settings.activities.index'),
    this.toRoute('settings.activities.activity'),
    this.use('toLeft'),
    this.reverse('toRight')
  );

  this.transition(
    this.fromRoute('settings.members.index'),
    this.toRoute('settings.members.member'),
    this.use('toLeft'),
    this.reverse('toRight')
  );
};

// app/components/_debug/probe.gjs
import Component from '@glimmer/component';

export default class DebugProbe extends Component {
  constructor() {
    super(...arguments);
    console.group('DebugProbe');
    console.log('typeof @value:', typeof this.args.value);
    if (this.args.value) {
      const proto = Object.getPrototypeOf(this.args.value);
      console.log('value prototype name:', proto?.constructor?.name);
      console.log('has updateProfile method?', typeof this.args.value.updateProfile);
      console.log('route name (if available):', this.args.value?.routeName ?? this.args.value?.currentRouteName);
    }
    console.groupEnd();
  }

  <template></template>
}

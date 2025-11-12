// app/routes/user/profile.js
import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class UserProfileRoute extends Route {
  @service api;

  constructor() {
    super(...arguments);
    console.log('app/routes/user/profile', 'constructed');
  }


  async model() {
    return this.api.getProfile();
  }

  get fiii() {
    console.log('app/routes/user/profile', 'fiii called');
    return true;
  }

  setupController(controller, model) {
    console.log('app/routes/user/profile', 'setupController');
    super.setupController(controller, model);
    controller._ping = () => console.log('controller._ping called');
    controller._ping();
    controller.fiii();
  }
}

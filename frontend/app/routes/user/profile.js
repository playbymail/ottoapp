// app/routes/user/profile.js
import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class UserProfileRoute extends Route {
  @service api;

  async model() {
    return this.api.getProfile();
  }
}

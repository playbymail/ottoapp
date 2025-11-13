// app/routes/user/profile.js
import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class UserProfileRoute extends Route {
  @service session;
  @service store;

  async model() {
    const userId = this.session.data.authenticated.user.id;
    return this.store.findRecord('user', userId);
  }
}

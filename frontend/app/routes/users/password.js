import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class UsersPasswordRoute extends Route {
  @service session;

  model() {
    return {
      userId: this.session.currentUserId,
    };
  }
}

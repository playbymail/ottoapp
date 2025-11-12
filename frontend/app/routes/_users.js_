import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class UsersRoute extends Route {
  @service session;
  @service router;

  beforeModel(transition) {
    // Require authentication first
    this.session.requireAuthentication(transition, 'login');

    // Then check if user has the "user" role
    if (!this.session.canAccessUserRoutes) {
      this.router.transitionTo('user.dashboard');
    }
  }
}

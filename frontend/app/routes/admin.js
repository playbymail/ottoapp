import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class AdminRoute extends Route {
  @service session;
  @service router;

  beforeModel(transition) {
    this.session.requireAuthentication(transition, 'login');

    // Check if user has the "admin" role
    if (!this.session.canAccessAdminRoutes) {
      this.router.transitionTo('user.dashboard');
    }
  }
}

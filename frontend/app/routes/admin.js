import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class AdminRoute extends Route {
  @service session;
  @service router;

  beforeModel(transition) {
    this.session.requireAuthentication(transition, 'login');

    // TODO: Add admin role check when user roles are implemented
    // if (!this.session.currentUser.isAdmin) {
    //   this.router.transitionTo('user.dashboard');
    // }
  }
}

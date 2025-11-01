import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class SecureRoute extends Route {
  @service session;
  @service router;

  beforeModel() {
    console.log('SecureRoute:beforeModel', 'this.session.isAuthenticated', this.session.isAuthenticated);
    if (!this.session.isAuthenticated) this.router.transitionTo('login');
  }
}

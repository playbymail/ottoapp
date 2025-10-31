import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class SecureRoute extends Route {
  @service session;
  beforeModel() {
    if (!this.session.isAuthenticated) this.transitionTo('login');
  }
}

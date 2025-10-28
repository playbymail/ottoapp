import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class SecureRoute extends Route {
  @service session;
  beforeModel() {
    if (!this.session.isAuthenticated) this.transitionTo('login');
  }
}

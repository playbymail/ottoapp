import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class SecureRoute extends Route {
  @service session;
  @service router;

  beforeModel() {
    console.log('esa', 'app/routes/secure:beforeModel');
    console.log('esa', 'app/routes/secure:beforeModel', 'session.isAuthenticated', this.session.isAuthenticated);
    console.log('esa', 'app/routes/secure:beforeModel', 'session.data', this.session.data);
    if (!this.session.isAuthenticated) this.router.transitionTo('login');
  }
}

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class ApplicationRoute extends Route {
  @service session;

  async beforeModel() {
    console.log('esa', 'app/routes/application:beforeModel');
    await this.session.setup();
    console.log('esa', 'app/authenticators/server:restore', 'session.isAuthenticated', this.session.isAuthenticated);
  }
}

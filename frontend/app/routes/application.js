import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class ApplicationRoute extends Route {
  @service session;
  @service currentUser;

  async beforeModel() {
    await this.session.setup();
    if (this.session.isAuthenticated) {
      await this.currentUser.load();
    }
  }
}

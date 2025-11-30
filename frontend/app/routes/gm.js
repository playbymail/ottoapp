// app/routes/gm.js

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class GmRoute extends Route {
  @service session;
  @service router;

  async beforeModel() {
    if (!this.session.canAccessGMRoutes) { // send them away if they’re not a GM
      this.router.transitionTo('login');
    }
  }

  model() {
    return { // any global GM data if needed (current game list, etc.)
      gameId: "0301", // default to TN3.1
    };
  }
}

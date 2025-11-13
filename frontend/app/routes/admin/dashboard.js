// app/routes/admin/dashboard.js

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class AdminDashboardIndexRoute extends Route {
  @service store;

  async model() {
    return this.store.findAll('user');
  }
}

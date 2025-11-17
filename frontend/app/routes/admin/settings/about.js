// app/routes/admin/settings/about.js

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class AdminSettingsAboutRoute extends Route {
  @service store;

  async model() {
    return this.store.findRecord('version', '1');
  }
}

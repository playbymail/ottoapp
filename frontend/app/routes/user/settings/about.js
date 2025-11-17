// app/routes/user/settings/about.js

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class UserSettingsAboutRoute extends Route {
  @service store;

  async model() {
    return this.store.findRecord('version', '1');
  }
}

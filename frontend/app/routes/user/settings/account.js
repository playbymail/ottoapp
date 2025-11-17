// app/routes/user/settings/account.js

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class UserSettingsAccountRoute extends Route {
  @service store;

  async model() {
    return this.store.findRecord('user', 'me');
  }
}

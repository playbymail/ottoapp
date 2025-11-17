// app/routes/user/settings/index.js

import Route from '@ember/routing/route';

export default class UserSettingsIndexRoute extends Route {
  beforeModel() {
    this.transitionTo('user.settings.account');
  }
}

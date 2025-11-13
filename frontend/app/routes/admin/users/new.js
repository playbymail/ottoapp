// app/routes/admin/users/new.js
import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class AdminUsersNewRoute extends Route {
  @service store;

  async model() {
    return this.store.createRecord('user', {
      username: '',
      email: '',
      timezone: 'Europe/London',
    });
  }
}

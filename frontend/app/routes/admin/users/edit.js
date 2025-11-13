// app/routes/admin/users/edit.js
import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class AdminUsersEditRoute extends Route {
  @service store;

  async model(params) {
    return this.store.findRecord('user', params.user_id);
  }
}

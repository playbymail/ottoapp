import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class AdminUsersEditRoute extends Route {
  @service api;

  async model(params) {
    const response = await this.api.getUser(params.user_id);
    if (response.ok) {
      return response.json();
    }
    throw new Error('Failed to load user');
  }
}

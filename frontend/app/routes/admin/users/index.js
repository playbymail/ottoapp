import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class AdminUsersIndexRoute extends Route {
  @service api;

  async model() {
    const response = await this.api.getUsers();
    if (response.ok) {
      return response.json();
    }
    throw new Error('Failed to load users');
  }
}

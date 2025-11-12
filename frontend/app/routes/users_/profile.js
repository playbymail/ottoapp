import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class UsersProfileRoute extends Route {
  @service api;

  async model() {
    const response = await this.api.getUserProfile();
    if (response.ok) {
      return response.json();
    }
    throw new Error('Failed to load user profile');
  }
}

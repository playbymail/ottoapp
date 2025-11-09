import AuthenticatedRoute from './authenticated';
import { service } from '@ember/service';

export default class ProfileRoute extends AuthenticatedRoute {
  @service store;

  async model() {
    // Fetch profile data from /api/profile
    const response = await fetch('/api/profile', {
      credentials: 'same-origin',
    });

    if (!response.ok) {
      throw new Error('Failed to fetch profile');
    }

    return response.json();
  }
}

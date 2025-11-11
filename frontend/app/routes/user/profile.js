import Route from '@ember/routing/route';

export default class UserProfileRoute extends Route {
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

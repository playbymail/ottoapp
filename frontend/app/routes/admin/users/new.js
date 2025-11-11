import Route from '@ember/routing/route';

export default class AdminUsersNewRoute extends Route {
  model() {
    return {
      username: '',
      email: '',
      password: '',
      timezone: 'UTC',
      roles: ['user'],
    };
  }
}

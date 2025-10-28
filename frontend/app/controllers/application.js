import Controller from '@ember/controller';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';

export default class ApplicationController extends Controller {
  @service session;
  @service currentUser;
  @service router;

  constructor() {
    super(...arguments);
    // Load user data once ESA restores a valid session
    this.session.on('authenticationSucceeded', () => this.currentUser.load());
    this.session.on('invalidationSucceeded', () => (this.currentUser.user = null));
  }

  @action async logout() {
    await this.session.invalidate();
    this.router.transitionTo('login');
  }
}

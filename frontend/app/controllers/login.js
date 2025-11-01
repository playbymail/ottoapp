// mdhender: not used?

import Controller from '@ember/controller';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

export default class LoginController extends Controller {
  @service session;
  @service router;

  @tracked username = '';
  @tracked password = '';
  @tracked error = null;

  @action updateUsername(e) { this.username = e.target.value; }
  @action updatePassword(e) { this.password = e.target.value; }

  @action async submit(e) {
    console.log('submit', this.username, this.password, this.session.isAuthenticated);
    e.preventDefault();
    this.error = null;
    try {
      await this.session.authenticate('authenticator:cookie', {
        username: this.username, password: this.password,
      });
    } catch {
      this.error = 'Invalid username or password.';
    }
    if (this.session.isAuthenticated) {
      this.router.transitionTo('secure');
    }
  }
}

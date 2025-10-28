import Controller from '@ember/controller';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

export default class LoginController extends Controller {
  @service session;

  @tracked username = '';
  @tracked password = '';
  @tracked error = null;

  @action updateUsername(e) { this.username = e.target.value; }
  @action updatePassword(e) { this.password = e.target.value; }

  @action async submit(e) {
    e.preventDefault();
    this.error = null;
    try {
      await this.session.authenticate('authenticator:cookie', {
        username: this.username, password: this.password,
      });
      this.transitionToRoute('secure');
    } catch {
      this.error = 'Invalid username or password.';
    }
  }
}

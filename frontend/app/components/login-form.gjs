import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { on } from '@ember/modifier';

export default class LoginFormComponent extends Component {
  @service session;
  @service router;

  @tracked username = '';
  @tracked password = '';
  @tracked error = null;

  @action updateUsername(e) { this.username = e.target.value; }
  @action updatePassword(e) { this.password = e.target.value; }

  // submit invokes session.authenticate with the user's credentials.
  // if the authentication succeeds, we route to the secure page.
  @action async submit(e) {
    e.preventDefault();
    this.error = null;
    try {
      await this.session.authenticate('authenticator:cookie', {
        username: this.username, password: this.password,
      });
      console.log('login-form', 'submit', this.session.isAuthenticated);
      this.router.transitionTo('secure');
    } catch {
      this.error = 'Invalid username or password.';
    }
  }

  <template>
    <h1>Sign in</h1>

    {{#if this.error}}
      <p role="alert">{{this.error}}</p>
    {{/if}}

    <form {{on "submit" this.submit}}>
      <label>
        Username
        <input type="text" value={{this.username}} {{on "input" this.updateUsername}} />
      </label>

      <label>
        Password
        <input type="password" value={{this.password}} {{on "input" this.updatePassword}} />
      </label>

      <button type="submit">Sign in</button>
    </form>
  </template>
}

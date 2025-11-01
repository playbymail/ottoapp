import Component from '@glimmer/component';
import {service} from '@ember/service';
import {tracked} from '@glimmer/tracking';
import {action} from '@ember/object';
import {on} from '@ember/modifier';

export default class LoginFormComponent extends Component {
  @service session;
  @service router;

  @tracked username = 'catbird';
  @tracked password = 'admin';
  @tracked error = null;
  @tracked avoidAutoFill = true;
  @tracked passwordType = 'text';

  @action updateUsername(e) {
    this.username = e.target.value;
  }

  @action updatePassword(e) {
    this.password = e.target.value;
  }

  // submit invokes session.authenticate with the user's credentials.
  // if the authentication succeeds, we route to the secure page.
  @action async submit(e) {
    e.preventDefault();
    this.error = null;
    try {
      await this.session.authenticate('authenticator:server', {
        username: this.username, password: this.password,
      });
      console.log('login-form', 'submit', this.session.isAuthenticated);
      this.router.transitionTo('secure');
    } catch {
      this.error = 'Invalid username or password.';
    }
  }

  <template>
    {{#if this.error}}
      <p role="alert">{{this.error}}</p>
    {{/if}}

    <div class="flex-grow flex items-center justify-center">
      <div class="bg-gray-800 p-8 rounded-lg shadow-lg w-96">
        <h1 class="text-2xl font-bold mb-6 text-center">OttoApp Login</h1>
        <form {{on "submit" this.submit}}
          autocomplete="off" data-1p-ignore data-lpignore="true">
          <div class="mb-4">
            <label for="username" class="block text-sm font-medium mb-2">Username</label>
            <input type="text" id="username" name="username" required
                   value={{this.username}} {{on "input" this.updateUsername}}
                   autocomplete="off" data-1p-ignore data-lpignore="true"
                   class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded focus:outline-none focus:border-blue-500">
          </div>
          <div class="mb-6">
            <label for="password" class="block text-sm font-medium mb-2">Password</label>
            <input type="text" id="password" name="password" required
                   value={{this.password}} {{on "input" this.updatePassword}}
                   autocomplete="off" data-1p-ignore data-lpignore="true"
                   class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded focus:outline-none focus:border-blue-500">
          </div>
          <button type="submit"
                  class="w-full bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 rounded transition">
            Login
          </button>
        </form>
      </div>
    </div>
  </template>
}

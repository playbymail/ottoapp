import Component from "@glimmer/component";
import {service} from "@ember/service";
import { on } from '@ember/modifier';
import { action } from "@ember/object";
import { LinkTo } from '@ember/routing';

export default class ApplicationHeaderComponent extends Component {
  @service session;
  @service router; // so we can redirect after logout

  get currentUser() {
    console.log('esa', 'app/components/application-header:getCurrentUser');
    console.log('esa', 'app/components/application-header:getCurrentUser', 'session.data', this.session.data);
    console.log('esa', 'app/components/application-header:getCurrentUser', 'session.data.authenticated', this.session.data.authenticated);
    console.log('esa', 'app/components/application-header:getCurrentUser', 'session.data.authenticated.user', this.session.data.authenticated.user);
    return this.session.data.authenticated.user;
  }
  get csrf() {
    return this.session.data.authenticated.csrf;
  }
  @action
  async logout() {
    // in dev, we can hit timing issues with the response, so ignore errors and assume the library did its job.
    await this.session.invalidate().catch(() => {}); // ← this calls app/authenticators/server.js → invalidate()
    this.router.transitionTo("/");  // force the browser back to the login page
  }

  <template>
    <header class="p-4 bg-gray-100 border-b flex justify-between">
      <h1 class="font-bold text-lg">Frontend Demo</h1>
      {{#if this.session.isAuthenticated}}
        <div class="flex items-center space-x-2">
          {{#if this.currentUser}}
            <span>👋 {{this.currentUser.username}}</span>
          {{/if}}
        </div>
        <button type="button" {{on "click" this.logout}}>
          Logout
        </button>
      {{else}}
        <LinkTo @route="login" class="underline text-blue-600">
          Login
        </LinkTo>
      {{/if}}
    </header>
    <div>
      <p>this.session.isAuthenticated: "{{this.session.isAuthenticated}}"</p>
      <p>this.session.data: "{{this.session.data}}"</p>
      <p>this.currentUser: "{{this.currentUser}}"</p>
      <p>this.currentUser: "{{this.currentUser.username}}"</p>
    </div>
  </template>
}

/*
          <button type="button" {{on "click" @onLogout}}>
            Logout
          </button>

 */

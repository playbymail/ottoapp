// Copyright (c) 2025 Michael D Henderson. All rights reserved.

import Component from '@glimmer/component';
import {service} from '@ember/service';
import {tracked} from '@glimmer/tracking';
import {action} from '@ember/object';
import {on} from '@ember/modifier';

// https://tailwindcss.com/plus/ui-blocks/application-ui/forms/sign-in-forms#simple
// Requires a TailwindCSS Plus license.

import { LinkTo } from '@ember/routing';

export default class LoginForm extends Component {
  @service session;
  @service router;

  @tracked email = '';
  @tracked password = '';
  @tracked error = null;

  @action updateEmail(e) {
    this.email = e.target.value;
  }

  @action updatePassword(e) {
    this.password = e.target.value;
  }

  // submit invokes session.authenticate with the user's credentials.
  // if the authentication succeeds, we route to the dashboard.
  @action async submit(e) {
    e.preventDefault();
    this.error = null;
    try {
      await this.session.authenticate('authenticator:server', {
        email: this.email, password: this.password,
      });
      console.log('login-form', 'submit', this.session.isAuthenticated);
      this.router.transitionTo('user.dashboard');
    } catch {
      this.error = 'Invalid credentials.';
    }
  }

  <template>
    <div class="mx-auto max-w-7xl sm:px-6 lg:px-8">
      <div
        class="flex min-h-full flex-col justify-center px-6 py-12 lg:px-8 bg-center bg-no-repeat bg-cover"
        style="background-image: url('/img/hero-bg-washed.jpg');">
        {{#if this.error}}
          <p role="alert">{{this.error}}</p>
        {{/if}}

        <div class="flex min-h-full flex-col justify-center px-6 py-12 lg:px-8">
          <div class="sm:mx-auto sm:w-full sm:max-w-sm">
            <img src="/img/logo-light.svg" alt="OttoApp" class="mx-auto h-10 w-auto dark:hidden" />
            <img src="/img/logo-dark.svg" alt="OttoApp" class="mx-auto h-10 w-auto not-dark:hidden" />
            <h2 class="mt-10 text-center text-2xl/9 font-bold tracking-tight text-gray-900">Sign in to your account</h2>
          </div>

          <div class="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
            <form {{on "submit" this.submit}} class="space-y-6">
              <div>
                <label for="email" class="block text-sm/6 font-medium text-gray-900">Email</label>
                <div class="mt-2">
                  <input id="email" type="text" name="email" required autocomplete="email"
                         value={{this.email}} {{on "input" this.updateEmail}}
                         class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6" />
                </div>
              </div>

              <div>
                <div class="flex items-center justify-between">
                  <label for="password" class="block text-sm/6 font-medium text-gray-900">Password</label>
                  <div class="text-sm">
                    <LinkTo @route="login" class="font-semibold text-indigo-600 hover:text-indigo-500">
                      Forgot password?
                    </LinkTo>
                  </div>
                </div>
                <div class="mt-2">
                  <input id="password" type="password" name="password" required autocomplete="current-password"
                         value={{this.password}} {{on "input" this.updatePassword}}
                         class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:text-sm/6" />
                </div>
              </div>

              <div>
                <button type="submit" class="flex w-full justify-center rounded-md bg-indigo-600 px-3 py-1.5 text-sm/6 font-semibold text-white shadow-xs hover:bg-indigo-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">Sign in</button>
              </div>
            </form>

            <p class="mt-10 text-center text-sm/6 text-gray-500">
              Need an account?
              <a href="https://discord.gg/xR9HVYGFXs" target="_blank" class="font-semibold text-indigo-600 hover:text-indigo-500" rel="noopener noreferrer">Request one here.</a>
            </p>
          </div>
        </div>
      </div>
    </div>
  </template>
}

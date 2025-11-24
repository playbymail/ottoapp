// Copyright (c) 2025 Michael D Henderson. All rights reserved.

import Component from '@glimmer/component';
import {service} from "@ember/service";
import { on } from '@ember/modifier';
import { action } from "@ember/object";
import not from 'frontend/helpers/not';

// https://tailwindcss.com/plus/ui-blocks/application-ui/application-shells/sidebar#sidebar-with-header
// Requires a TailwindCSS Plus license.

import { LinkTo } from '@ember/routing';

function minimumDelay(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export default class ProfileDropdown extends Component {
  @service session;
  @service router; // so we can redirect after logout

  get currentUser() {
    if (!this.session.isAuthenticated) {
      return {
        handle: 'Guest',
      }
    }
    return this.session.data.authenticated.user;
  }

  get csrf() {
    return this.session.data.authenticated.csrf;
  }

  @action async logout() {
    // note: ESA forces a route change after logout, so the route in
    // the LinkTo that calls us gets ignored.

    // create a promise to call app/authenticators/server.js → invalidate()
    let invalidatePromise = this.session.invalidate().catch(() => {
      // Ignore errors during development. We can hit timing issues
      // with the response, so we must ignore errors and assume the
      // library did its job.
    });
    // run the promise in parallel with our minimum delay, ensuring at
    // least 250ms passes so that any background workers can finish.
    await Promise.all([invalidatePromise, minimumDelay(250)]);
  }

  <template>
    {{!-- Profile dropdown --}}
    <el-dropdown class="relative">
      <button class="relative flex items-center" type="button">
        <span class="absolute -inset-1.5"></span>
        <span class="sr-only">Open user menu</span>
        <img src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80" alt="" class="size-8 rounded-full bg-gray-50 outline -outline-offset-1 outline-black/5 dark:bg-gray-800 dark:outline-white/10" />
        <span class="hidden lg:flex lg:items-center">
          <span aria-hidden="true" class="ml-4 text-sm/6 font-semibold text-gray-900 dark:text-white">
            {{#if this.session.isAuthenticated}}{{this.session.getHandle}}{{else}}guest{{/if}}
          </span>
          <svg viewBox="0 0 20 20" fill="currentColor" data-slot="icon" aria-hidden="true" class="ml-2 size-5 text-gray-400 dark:text-gray-500">
            <path d="M5.22 8.22a.75.75 0 0 1 1.06 0L10 11.94l3.72-3.72a.75.75 0 1 1 1.06 1.06l-4.25 4.25a.75.75 0 0 1-1.06 0L5.22 9.28a.75.75 0 0 1 0-1.06Z" clip-rule="evenodd" fill-rule="evenodd" />
          </svg>
        </span>
      </button>
      <el-menu anchor="bottom end" popover class="w-32 origin-top-right rounded-md bg-white py-2 shadow-lg outline-1 outline-gray-900/5 transition transition-discrete [--anchor-gap:--spacing(2.5)] data-closed:scale-95 data-closed:transform data-closed:opacity-0 data-enter:duration-100 data-enter:ease-out data-leave:duration-75 data-leave:ease-in dark:bg-gray-800 dark:shadow-none dark:-outline-offset-1 dark:outline-white/10">
        {{#if this.session.isAuthenticated }}
          <LinkTo @route="user.settings.account" class="block px-3 py-1 text-sm/6 text-gray-900 focus:bg-gray-50 focus:outline-hidden dark:text-white dark:focus:bg-white/5">
            My profile
          </LinkTo>
          <LinkTo @route="login" {{on 'click' this.logout}}
                  class="block px-3 py-1 text-sm/6 text-gray-900 focus:bg-gray-50 focus:outline-hidden dark:text-white dark:focus:bg-white/5">
            Sign out
          </LinkTo>
        {{else}}
          <LinkTo @route="login" class="block px-3 py-1 text-sm/6 text-gray-900 focus:bg-gray-50 focus:outline-hidden dark:text-white dark:focus:bg-white/5">
            Sign in
          </LinkTo>
        {{/if}}
      </el-menu>
    </el-dropdown>
  </template>
}

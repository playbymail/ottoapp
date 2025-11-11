// Copyright (c) 2025 Michael D Henderson. All rights reserved.
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { on } from '@ember/modifier';

// Requires a TailwindCSS Plus license.

import {LinkTo} from '@ember/routing';

export default class Password extends Component {
  @service api;
  @service session;

  <template>
    <div class="bg-white shadow sm:rounded-lg">
      <div class="px-4 py-5 sm:p-6">
        <h3 class="text-lg font-medium leading-6 text-gray-900">
          Change Password
        </h3>
        <div class="mt-5">
          {{#if this.errorMessage}}
            <div class="rounded-md bg-red-50 p-4 mb-4">
              <p class="text-sm text-red-700">{{@controller.errorMessage}}</p>
            </div>
          {{/if}}
          {{#if this.successMessage}}
            <div class="rounded-md bg-green-50 p-4 mb-4">
              <p class="text-sm text-green-700">{{@controller.successMessage}}</p>
            </div>
          {{/if}}

          <form {{on "submit" this.changePassword}}>
            <div class="space-y-4">
              <div>
                <label for="currentPassword" class="block text-sm font-medium text-gray-700">
                  Current Password
                </label>
                <input
                  type="password"
                  name="currentPassword"
                  id="currentPassword"
                  required
                  class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                />
              </div>

              <div>
                <label for="newPassword" class="block text-sm font-medium text-gray-700">
                  New Password
                </label>
                <input
                  type="password"
                  name="newPassword"
                  id="newPassword"
                  required
                  class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                />
              </div>

              <div>
                <label for="confirmPassword" class="block text-sm font-medium text-gray-700">
                  Confirm New Password
                </label>
                <input
                  type="password"
                  name="confirmPassword"
                  id="confirmPassword"
                  required
                  class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                />
              </div>

              <div class="flex justify-end">
                <button
                  type="submit"
                  disabled={{this.isSaving}}
                  class="inline-flex justify-center rounded-md border border-transparent bg-indigo-600 py-2 px-4 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 disabled:opacity-50"
                >
                  {{if this.isSaving "Changing..." "Change Password"}}
                </button>
              </div>
            </div>
          </form>
        </div>
      </div>
    </div>
  </template>
}

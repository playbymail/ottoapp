// Copyright (c) 2025 Michael D Henderson. All rights reserved.

import Component from '@glimmer/component';
import {service} from "@ember/service";
import { on } from '@ember/modifier';
import { action } from "@ember/object";

// Requires a TailwindCSS Plus license.

import { LinkTo } from '@ember/routing';

export default class New extends Component {
  <template>
    <div class="bg-white shadow sm:rounded-lg">
      <div class="px-4 py-5 sm:p-6">
        <h3 class="text-lg font-medium leading-6 text-gray-900">
          Create New User
        </h3>
        <div class="mt-5">
          {{#if @controller.errorMessage}}
            <div class="rounded-md bg-red-50 p-4 mb-4">
              <p class="text-sm text-red-700">{{@controller.errorMessage}}</p>
            </div>
          {{/if}}

          <form {{on "submit" @controller.createUser}}>
            <div class="space-y-4">
              <div>
                <label for="username" class="block text-sm font-medium text-gray-700">
                  Username
                </label>
                <input
                  type="text"
                  name="username"
                  id="username"
                  required
                  class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                />
              </div>

              <div>
                <label for="email" class="block text-sm font-medium text-gray-700">
                  Email
                </label>
                <input
                  type="email"
                  name="email"
                  id="email"
                  required
                  class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                />
              </div>

              <div>
                <label for="password" class="block text-sm font-medium text-gray-700">
                  Password (leave empty to auto-generate)
                </label>
                <input
                  type="password"
                  name="password"
                  id="password"
                  class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                />
              </div>

              <div>
                <label for="timezone" class="block text-sm font-medium text-gray-700">
                  Timezone
                </label>
                <input
                  type="text"
                  name="timezone"
                  id="timezone"
                  value="UTC"
                  required
                  class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                />
              </div>

              <div class="flex justify-end gap-3">
                <button
                  type="button"
                  {{on "click" @controller.cancel}}
                  class="inline-flex justify-center rounded-md border border-gray-300 bg-white py-2 px-4 text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={{@controller.isSaving}}
                  class="inline-flex justify-center rounded-md border border-transparent bg-indigo-600 py-2 px-4 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 disabled:opacity-50"
                >
                  {{if @controller.isSaving "Creating..." "Create User"}}
                </button>
              </div>
            </div>
          </form>
        </div>
      </div>
    </div>
  </template>
}

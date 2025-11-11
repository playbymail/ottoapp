// Copyright (c) 2025 Michael D Henderson. All rights reserved.

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { on } from '@ember/modifier';

// Requires a TailwindCSS Plus license.

import {LinkTo} from '@ember/routing';

export default class Profile extends Component {
  @service api;
  @service session;

  // model / state variables?
  @tracked username;
  @tracked email = "";
  @tracked timezone = "";


  // save/submit variables
  @tracked errorMessages = [];
  @tracked successMessage = "";

  constructor() {
    super(...arguments);

    console.log('users.profile.args.model', this.args.model);
    // transfer incoming model data to local tracked variables
    this.email = this.args.model?.email || "";
    this.timezone = this.args.model?.timezone || "";
    this.username = this.args.model?.username || "";
  }

  get hasChanges() {
    return this.email !== this.profile.email || this.timezone !== this.profile.timezone || this.username !== this.username;
  }

  @action
  async updateProfile(event) {
    // agent generated an @controller.updateProfile that we're trying to replace
    // other code uses an @action async save
    console.log('implement me');
  }

  <template>
    <div class="bg-white shadow sm:rounded-lg">
      <div class="px-4 py-5 sm:p-6">
        <h3 class="text-lg font-medium leading-6 text-gray-900">
          Edit Profile
        </h3>
        <div class="mt-5">
          {{#if this.errorMessage}}
            <div class="rounded-md bg-red-50 p-4 mb-4">
              <p class="text-sm text-red-700">{{this.errorMessage}}</p>
            </div>
          {{/if}}
          {{#if this.successMessage}}
            <div class="rounded-md bg-green-50 p-4 mb-4">
              <p class="text-sm text-green-700">{{this.successMessage}}</p>
            </div>
          {{/if}}

          <form {{on "submit" this.updateProfile}}>
            <div class="space-y-4">
              <div>
                <label for="username" class="block text-sm font-medium text-gray-700">
                  Username
                </label>
                <input
                  type="text"
                  name="username"
                  id="username"
                  value={{@model.username}}
                  disabled
                  class="mt-1 block w-full rounded-md border-gray-300 bg-gray-100 shadow-sm sm:text-sm"
                />
                <p class="mt-1 text-sm text-gray-500">Username cannot be changed</p>
              </div>

              <div>
                <label for="email" class="block text-sm font-medium text-gray-700">
                  Email
                </label>
                <input
                  type="email"
                  name="email"
                  id="email"
                  value={{@model.email}}
                  required
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
                  value={{@model.timezone}}
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
                  {{if this.isSaving "Saving..." "Save Changes"}}
                </button>
              </div>
            </div>
          </form>
        </div>
      </div>
    </div>
  </template>
}

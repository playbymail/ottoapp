// Copyright (c) 2025 Michael D Henderson. All rights reserved.

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { on } from '@ember/modifier';

import TimezonePicker  from 'frontend/components/timezone-picker';


export default class Profile extends Component {
  @service session;

  @tracked email = "";
  @tracked timezone = "";
  @tracked isSaving = false;
  @tracked errorMessages = [];
  @tracked successMessage = "";

  constructor() {
    super(...arguments);
    this.email = this.args.model?.email || "";
    this.timezone = this.args.model?.timezone || "";
  }

  get profile() {
    return this.args.model || {};
  }

  get hasChanges() {
    return this.email !== this.profile.email || this.timezone !== this.profile.timezone;
  }

  get username() {
    return this.profile.username || "";
  }

  @action
  updateEmail(event) {
    this.email = event.target.value;
    this.errorMessages = [];
    this.successMessage = "";
  }

  @action
  updateTimezone(label) {
    this.timezone = label;
    this.errorMessages = [];
    this.successMessage = "";
  }

  @action
  dismissSuccess() {
    this.successMessage = "";
  }

  @action
  async save(event) {
    event.preventDefault();

    if (!this.hasChanges) {
      return;
    }

    this.isSaving = true;
    this.errorMessages = [];

    try {
      const response = await fetch('/api/profile', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-CSRF-Token': this.session.data.authenticated.csrf,
        },
        credentials: 'same-origin',
        body: JSON.stringify({
          email: this.email,
          timezone: this.timezone,
        }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        this.errorMessages = errorData.errors || ['Failed to update profile'];
        return;
      }

      const data = await response.json();
      // Update the model with all fields from the response
      this.args.model.username = data.username;
      this.args.model.email = data.email;
      this.args.model.timezone = data.timezone;
      this.email = data.email;
      this.timezone = data.timezone;
      this.successMessage = "Profile updated successfully";

      // Auto-dismiss success message after 3 seconds
      setTimeout(() => {
        this.successMessage = "";
      }, 3000);
    } catch (error) {
      this.errorMessages = [error.message || 'An unexpected error occurred'];
    } finally {
      this.isSaving = false;
    }
  }

  @action
  cancel() {
    this.email = this.profile.email;
    this.timezone = this.profile.timezone;
    this.errorMessages = [];
    this.successMessage = "";
  }

  <template>
    <form {{on "submit" this.save}} class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <div class="space-y-12 sm:space-y-16">
        <div>
          <h2 class="text-base/7 font-semibold text-gray-900">Profile</h2>
          <p class="mt-1 max-w-2xl text-sm/6 text-gray-600">
            The information in this section will be displayed publicly, so be careful what you share.
          </p>

          <div class="mt-10 space-y-8 border-b border-gray-900/10 pb-12 sm:space-y-0 sm:divide-y sm:divide-gray-900/10 sm:border-t sm:border-t-gray-900/10 sm:pb-0">
            <div class="sm:grid sm:grid-cols-3 sm:items-start sm:gap-4 sm:py-6">
              <label for="username" class="block text-sm/6 font-medium text-gray-900 sm:pt-1.5">Username</label>
              <div class="mt-2 sm:col-span-2 sm:mt-0">
                <input
                  id="username"
                  type="text"
                  name="username"
                  value={{this.username}}
                  disabled
                  class="block w-full rounded-md bg-gray-50 px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 sm:max-w-md sm:text-sm/6"
                />
              </div>
            </div>
          </div>
        </div>

        <div>
          <h2 class="text-base/7 font-semibold text-gray-900">Personal Information</h2>
          <p class="mt-1 max-w-2xl text-sm/6 text-gray-600">
            The information is this section is not shared.
            Use a permanent address where you can receive mail.
          </p>

          <div class="mt-10 space-y-8 border-b border-gray-900/10 pb-12 sm:space-y-0 sm:divide-y sm:divide-gray-900/10 sm:border-t sm:border-t-gray-900/10 sm:pb-0">
            <div class="sm:grid sm:grid-cols-3 sm:items-start sm:gap-4 sm:py-6">
              <label for="email" class="block text-sm/6 font-medium text-gray-900 sm:pt-1.5">Email address</label>
              <div class="mt-2 sm:col-span-2 sm:mt-0">
                <input
                  id="email"
                  type="email"
                  name="email"
                  value={{this.email}}
                  {{on "input" this.updateEmail}}
                  autocomplete="email"
                  class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:max-w-md sm:text-sm/6"
                />
              </div>
            </div>

            <div class="sm:grid sm:grid-cols-3 sm:items-start sm:gap-4 sm:py-6">
              <label for="timezone" class="block text-sm/6 font-medium text-gray-900 sm:pt-1.5">Timezone</label>
              <div class="mt-2 sm:col-span-2 sm:mt-0">
                <TimezonePicker
                  @value={{this.timezone}}
                  @onChange={{this.updateTimezone}}
                />
              </div>
            </div>
          </div>
        </div>
      </div>

      {{#if this.successMessage}}
        <div class="mt-6 rounded-md bg-green-50 p-4">
          <div class="flex">
            <div class="shrink-0">
              <svg viewBox="0 0 20 20" fill="currentColor" data-slot="icon" aria-hidden="true" class="size-5 text-green-400">
                <path d="M10 18a8 8 0 1 0 0-16 8 8 0 0 0 0 16Zm3.857-9.809a.75.75 0 0 0-1.214-.882l-3.483 4.79-1.88-1.88a.75.75 0 1 0-1.06 1.061l2.5 2.5a.75.75 0 0 0 1.137-.089l4-5.5Z" clip-rule="evenodd" fill-rule="evenodd" />
              </svg>
            </div>
            <div class="ml-3">
              <p class="text-sm font-medium text-green-800">{{this.successMessage}}</p>
            </div>
            <div class="ml-auto pl-3">
              <div class="-mx-1.5 -my-1.5">
                <button type="button" {{on "click" this.dismissSuccess}} class="inline-flex rounded-md bg-green-50 p-1.5 text-green-500 hover:bg-green-100 focus-visible:ring-2 focus-visible:ring-green-600 focus-visible:ring-offset-2 focus-visible:ring-offset-green-50 focus-visible:outline-hidden">
                  <span class="sr-only">Dismiss</span>
                  <svg viewBox="0 0 20 20" fill="currentColor" data-slot="icon" aria-hidden="true" class="size-5">
                    <path d="M6.28 5.22a.75.75 0 0 0-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 1 0 1.06 1.06L10 11.06l3.72 3.72a.75.75 0 1 0 1.06-1.06L11.06 10l3.72-3.72a.75.75 0 0 0-1.06-1.06L10 8.94 6.28 5.22Z" />
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>
      {{/if}}

      {{#if this.errorMessages.length}}
        <div class="mt-6 rounded-md bg-red-50 p-4">
          <div class="flex">
            <div class="shrink-0">
              <svg viewBox="0 0 20 20" fill="currentColor" data-slot="icon" aria-hidden="true" class="size-5 text-red-400">
                <path d="M10 18a8 8 0 1 0 0-16 8 8 0 0 0 0 16ZM8.28 7.22a.75.75 0 0 0-1.06 1.06L8.94 10l-1.72 1.72a.75.75 0 1 0 1.06 1.06L10 11.06l1.72 1.72a.75.75 0 1 0 1.06-1.06L11.06 10l1.72-1.72a.75.75 0 0 0-1.06-1.06L10 8.94 8.28 7.22Z" clip-rule="evenodd" fill-rule="evenodd" />
              </svg>
            </div>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-red-800">The update failed</h3>
              <div class="mt-2 text-sm text-red-700">
                <ul role="list" class="list-disc space-y-1 pl-5">
                  {{#each this.errorMessages as |msg|}}
                    <li>{{msg}}</li>
                  {{/each}}
                </ul>
              </div>
            </div>
          </div>
        </div>
      {{/if}}

      <div class="mt-6 flex items-center justify-end gap-x-6">
        <button
          type="button"
          {{on "click" this.cancel}}
          disabled={{unless this.hasChanges "disabled"}}
          class="text-sm/6 font-semibold text-gray-900 disabled:text-gray-400 disabled:cursor-not-allowed"
        >
          Cancel
        </button>
        <button
          type="submit"
          disabled={{if this.isSaving "disabled" (unless this.hasChanges "disabled")}}
          class="inline-flex justify-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-xs hover:bg-indigo-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 disabled:bg-gray-400 disabled:cursor-not-allowed"
        >
          {{if this.isSaving "Saving..." "Save"}}
        </button>
      </div>
    </form>
  </template>
}

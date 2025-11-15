// app/components/user/form.gjs
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { on } from '@ember/modifier';

import TimezonePicker from 'frontend/components/timezone-picker';

export default class UserForm extends Component {
  @tracked email = this.args.model?.email ?? '';
  @tracked timezone = this.args.model?.timezone ?? '';
  @tracked username = this.args.model?.username ?? '';
  @tracked permissions = this.args.model?.permissions ?? {};

  get hasChanges() {
    const emailChanged = this.email !== (this.args.model?.email ?? '');
    const timezoneChanged = this.timezone !== (this.args.model?.timezone ?? '');
    const usernameChanged = this.args.canEditUsername && this.username !== (this.args.model?.username ?? '');

    return emailChanged || timezoneChanged || usernameChanged;
  }

  get isCancelDisabled() {
    return !this.hasChanges;
  }

  get isSaveDisabled() {
    return this.args.isSaving || !this.hasChanges;
  }

  get errorList() {
    const e = this.args.errorMessage;
    if (!e) return [];
    return Array.isArray(e) ? e : [e];
  }

  get profile() {
    return this.args.model ?? {};
  }

  @action updateEmail(e) {
    this.email = e.target.value;
  }

  @action updateUsername(e) {
    this.username = e.target.value;
  }

  @action updateTimezone(label) {
    this.timezone = label;
  }

  @action async save(e) {
    e?.preventDefault();
    e?.stopPropagation();

    if (!this.hasChanges) return;

    // Push values into the model
    this.args.model.email = this.email;
    this.args.model.timezone = this.timezone;
    if (this.args.canEditUsername) {
      this.args.model.username = this.username;
    }

    await this.args.onSave?.();
  }

  @action cancel() {
    this.email = this.args.model?.email ?? '';
    this.timezone = this.args.model?.timezone ?? '';
    this.username = this.args.model?.username ?? '';
    this.args.onCancel?.();
  }

  <template>
    <form {{on "submit" this.save}} class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <div class="space-y-12 sm:space-y-16">
        <div>
          <h2 class="text-base/7 font-semibold text-gray-900">Profile</h2>
          <p class="mt-1 max-w-2xl text-sm/6 text-gray-600">
            {{#if @canEditUsername}}
              Edit user account information.
            {{else}}
              The information in this section will be displayed publicly, so be careful what you share.
            {{/if}}
          </p>

          <div class="mt-10 space-y-8 border-b border-gray-900/10 pb-12 sm:space-y-0 sm:divide-y sm:divide-gray-900/10 sm:border-t sm:border-t-gray-900/10 sm:pb-0">
            <div class="sm:grid sm:grid-cols-3 sm:items-start sm:gap-4 sm:py-6">
              <label for="handle" class="block text-sm/6 font-medium text-gray-900 sm:pt-1.5">Handle</label>
              <div class="mt-2 sm:col-span-2 sm:mt-0">
                <input
                  id="handle"
                  type="text"
                  name="handle"
                  value={{this.profile.handle}}
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
            The information in this section is not shared.
          </p>

          <div class="sm:grid sm:grid-cols-3 sm:items-start sm:gap-4 sm:py-6">
            <label for="username" class="block text-sm/6 font-medium text-gray-900 sm:pt-1.5">Username</label>
            <div class="mt-2 sm:col-span-2 sm:mt-0">
              {{#if @canEditUsername}}
                <input
                  id="username"
                  type="text"
                  name="username"
                  value={{this.username}}
                  {{on "input" this.updateUsername}}
                  class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600 sm:max-w-md sm:text-sm/6"
                />
              {{else}}
                <input
                  id="username"
                  type="text"
                  name="username"
                  value={{this.username}}
                  disabled
                  class="block w-full rounded-md bg-gray-50 px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 sm:max-w-md sm:text-sm/6"
                />
              {{/if}}
            </div>
          </div>

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

        {{#if this.permissions}}
          <div>
            <h2 class="text-base/7 font-semibold text-gray-900">OttoMap Permissions</h2>
            <p class="mt-1 max-w-2xl text-sm/6 text-gray-600">
              Only administrators may make updates.
            </p>

            <div class="mt-10 space-y-8 border-b border-gray-900/10 pb-12 sm:space-y-0 sm:divide-y sm:divide-gray-900/10 sm:border-t sm:border-t-gray-900/10 sm:pb-0">
              <div class="sm:grid sm:grid-cols-3 sm:items-start sm:gap-4 sm:py-6">
                <p class="text-sm/6 font-medium text-gray-900">Roles</p>
                <div class="mt-2 sm:col-span-2 sm:mt-0">
                  {{#if this.profile.roles}}
                    {{#each this.profile.roles as |role|}}
                      <span class="inline-flex items-center rounded-full bg-indigo-100 px-2.5 py-0.5 text-xs font-medium text-indigo-800 mr-1">
                        {{role}}
                      </span>
                    {{/each}}
                  {{else}}
                    <span class="text-sm text-gray-500">No roles assigned</span>
                  {{/if}}
                </div>
              </div>
            </div>
          </div>
        {{/if}}
      </div>

      {{!-- success from controller --}}
      {{#if @successMessage}}
        <div class="mt-6 rounded-md bg-green-50 p-4">
          <div class="flex">
            <div class="shrink-0">
              <svg viewBox="0 0 20 20" aria-hidden="true" class="size-5 text-green-400" fill="currentColor">
                <path d="M10 18a8 8 0 1 0 0-16 8 8 0 0 0 0 16Zm3.857-9.809a.75.75 0 0 0-1.214-.882l-3.483 4.79-1.88-1.88a.75.75 0 1 0-1.06 1.061l2.5 2.5a.75.75 0 0 0 1.137-.089l4-5.5Z" />
              </svg>
            </div>
            <div class="ml-3">
              <p class="text-sm font-medium text-green-800">
                {{@successMessage}}
              </p>
            </div>
          </div>
        </div>
      {{/if}}

      {{!-- errors from controller; accept string or array --}}
      {{#if this.errorList.length}}
        <div class="mt-6 rounded-md bg-red-50 p-4">
          <div class="flex">
            <div class="shrink-0">
              <svg viewBox="0 0 20 20" aria-hidden="true" class="size-5 text-red-400" fill="currentColor">
                <path d="M10 18a8 8 0 1 0 0-16 8 8 0 0 0 0 16ZM8.28 7.22a.75.75 0 0 0-1.06 1.06L8.94 10l-1.72 1.72a.75.75 0 1 0 1.06 1.06L10 11.06l1.72 1.72a.75.75 0 1 0 1.06-1.06L11.06 10l1.72-1.72a.75.75 0 0 0-1.06-1.06L10 8.94 8.28 7.22Z" />
              </svg>
            </div>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-red-800">The update failed</h3>
              <ul role="list" class="mt-2 list-disc space-y-1 pl-5 text-sm text-red-700">
                {{#each this.errorList as |msg|}}
                  <li>{{msg}}</li>
                {{/each}}
              </ul>
            </div>
          </div>
        </div>
      {{/if}}

      <div class="mt-6 flex items-center justify-end gap-x-6">
        <button
          type="button"
          {{on "click" this.cancel}}
          disabled={{this.isCancelDisabled}}
          class="text-sm/6 font-semibold text-gray-900 disabled:text-gray-400 disabled:cursor-not-allowed"
        >
          Cancel
        </button>
        <button
          type="submit"
          disabled={{this.isSaveDisabled}}
          class="inline-flex justify-center rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-xs hover:bg-indigo-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 disabled:bg-gray-400 disabled:cursor-not-allowed"
        >
          {{if @isSaving "Saving..." "Save"}}
        </button>
      </div>
    </form>
  </template>
}

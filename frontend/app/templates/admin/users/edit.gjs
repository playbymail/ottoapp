import RouteTemplate from 'ember-route-template';

export default RouteTemplate(
  <template>
    <div class="bg-white shadow sm:rounded-lg">
      <div class="px-4 py-5 sm:p-6">
        <h3 class="text-lg font-medium leading-6 text-gray-900">
          Edit User: {{@model.username}}
        </h3>
        <div class="mt-5">
          {{#if @controller.errorMessage}}
            <div class="rounded-md bg-red-50 p-4 mb-4">
              <p class="text-sm text-red-700">{{@controller.errorMessage}}</p>
            </div>
          {{/if}}
          {{#if @controller.successMessage}}
            <div class="rounded-md bg-green-50 p-4 mb-4">
              <p class="text-sm text-green-700">{{@controller.successMessage}}</p>
            </div>
          {{/if}}

          <form {{on "submit" @controller.updateUser}}>
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
                  disabled={{not @model.permissions.canEditUsername}}
                  required
                  class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm disabled:bg-gray-100"
                />
                {{#unless @model.permissions.canEditUsername}}
                  <p class="mt-1 text-sm text-gray-500">You cannot edit this user's username</p>
                {{/unless}}
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

              <div>
                <label class="block text-sm font-medium text-gray-700">
                  Roles
                </label>
                <div class="mt-2">
                  {{#each @model.roles as |role|}}
                    <span class="inline-flex items-center rounded-full bg-indigo-100 px-2.5 py-0.5 text-xs font-medium text-indigo-800 mr-1">
                      {{role}}
                    </span>
                  {{/each}}
                </div>
              </div>

              <div class="flex justify-between">
                <div class="flex gap-3">
                  {{#if @model.permissions.canResetPassword}}
                    <button
                      type="button"
                      {{on "click" @controller.resetPassword}}
                      disabled={{@controller.isSaving}}
                      class="inline-flex justify-center rounded-md border border-gray-300 bg-white py-2 px-4 text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 disabled:opacity-50"
                    >
                      Reset Password
                    </button>
                  {{/if}}
                </div>
                <div class="flex gap-3">
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
                    {{if @controller.isSaving "Saving..." "Save Changes"}}
                  </button>
                </div>
              </div>
            </div>
          </form>

          {{#if @controller.showResetPasswordDialog}}
            <div class="fixed inset-0 bg-gray-500 bg-opacity-75 flex items-center justify-center z-50">
              <div class="bg-white rounded-lg p-6 max-w-md">
                <h3 class="text-lg font-medium text-gray-900 mb-4">
                  Password Reset
                </h3>
                <p class="text-sm text-gray-700 mb-4">
                  The temporary password has been generated. Please provide this to the user:
                </p>
                <div class="bg-gray-100 rounded p-4 mb-4">
                  <code class="text-lg font-mono">{{@controller.tempPassword}}</code>
                </div>
                <div class="flex justify-end">
                  <button
                    type="button"
                    {{on "click" @controller.closeResetPasswordDialog}}
                    class="inline-flex justify-center rounded-md border border-transparent bg-indigo-600 py-2 px-4 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
                  >
                    Close
                  </button>
                </div>
              </div>
            </div>
          {{/if}}
        </div>
      </div>
    </div>
  </template>
);

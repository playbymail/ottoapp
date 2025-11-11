import RouteTemplate from 'ember-route-template';

export default RouteTemplate(
  <template>
    <div class="bg-white shadow sm:rounded-lg">
      <div class="px-4 py-5 sm:p-6">
        <h3 class="text-lg font-medium leading-6 text-gray-900">
          Edit Profile
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

          <form {{on "submit" @controller.updateProfile}}>
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
                  disabled={{@controller.isSaving}}
                  class="inline-flex justify-center rounded-md border border-transparent bg-indigo-600 py-2 px-4 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 disabled:opacity-50"
                >
                  {{if @controller.isSaving "Saving..." "Save Changes"}}
                </button>
              </div>
            </div>
          </form>
        </div>
      </div>
    </div>
  </template>
);

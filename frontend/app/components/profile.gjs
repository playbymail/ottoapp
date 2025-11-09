// Copyright (c) 2025 Michael D Henderson. All rights reserved.

import Component from "@glimmer/component";

export default class Profile extends Component {
  get profile() {
    return this.args.model || {};
  }

  <template>
    <form class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <div class="space-y-12 sm:space-y-16">
        <div>
          <h2 class="text-base/7 font-semibold text-gray-900">Profile</h2>
          <p class="mt-1 max-w-2xl text-sm/6 text-gray-600">This information will be displayed publicly so be careful what you share.</p>

          <div class="mt-10 space-y-8 border-b border-gray-900/10 pb-12 sm:space-y-0 sm:divide-y sm:divide-gray-900/10 sm:border-t sm:border-t-gray-900/10 sm:pb-0">
            <div class="sm:grid sm:grid-cols-3 sm:items-start sm:gap-4 sm:py-6">
              <label for="username" class="block text-sm/6 font-medium text-gray-900 sm:pt-1.5">Username</label>
              <div class="mt-2 sm:col-span-2 sm:mt-0">
                <input
                  id="username"
                  type="text"
                  name="username"
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
          <p class="mt-1 max-w-2xl text-sm/6 text-gray-600">Use a permanent address where you can receive mail.</p>

          <div class="mt-10 space-y-8 border-b border-gray-900/10 pb-12 sm:space-y-0 sm:divide-y sm:divide-gray-900/10 sm:border-t sm:border-t-gray-900/10 sm:pb-0">
            <div class="sm:grid sm:grid-cols-3 sm:items-start sm:gap-4 sm:py-6">
              <label for="email" class="block text-sm/6 font-medium text-gray-900 sm:pt-1.5">Email address</label>
              <div class="mt-2 sm:col-span-2 sm:mt-0">
                <input
                  id="email"
                  type="email"
                  name="email"
                  value={{this.profile.email}}
                  disabled
                  autocomplete="email"
                  class="block w-full rounded-md bg-gray-50 px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 sm:max-w-md sm:text-sm/6"
                />
              </div>
            </div>

            <div class="sm:grid sm:grid-cols-3 sm:items-start sm:gap-4 sm:py-6">
              <label for="timezone" class="block text-sm/6 font-medium text-gray-900 sm:pt-1.5">Timezone</label>
              <div class="mt-2 sm:col-span-2 sm:mt-0">
                <input
                  id="timezone"
                  type="text"
                  name="timezone"
                  value={{this.profile.timezone}}
                  disabled
                  class="block w-full rounded-md bg-gray-50 px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 sm:max-w-xs sm:text-sm/6"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </form>
  </template>
}

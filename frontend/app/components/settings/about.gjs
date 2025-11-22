// app/components/settings/about.gjs

import Component from "@glimmer/component";

export default class SettingsAboutTab extends Component {
  <template>
    <div class="grid max-w-7xl grid-cols-1 gap-x-8 gap-y-10 px-4 py-16 sm:px-6 md:grid-cols-3 lg:px-8">
      <div>
        <h2 class="text-base/7 font-semibold text-gray-900">Application Version</h2>
        <p class="mt-1 text-sm/6 text-gray-500">Version information for the browser client.</p>
      </div>

      <div class="md:col-span-2">
        <dl class="grid grid-cols-1 gap-x-6 gap-y-8 sm:max-w-xl sm:grid-cols-6">
          <div class="col-span-full">
            <dt class="block text-sm/6 font-medium text-gray-900">Version</dt>
            <dd class="mt-2 text-sm/6 text-gray-700">{{@version.short}}</dd>
          </div>
        </dl>
      </div>

      <div>
        <h2 class="text-base/7 font-semibold text-gray-900">Server Version</h2>
        <p class="mt-1 text-sm/6 text-gray-500">Version information for the API server.</p>
      </div>

      <div class="md:col-span-2">
        <dl class="grid grid-cols-1 gap-x-6 gap-y-8 sm:max-w-xl sm:grid-cols-6">
          <div class="col-span-full">
            <dt class="block text-sm/6 font-medium text-gray-900">Build</dt>
            <dd class="mt-2 text-sm/6 text-gray-700">{{@version.full}}</dd>
          </div>
        </dl>
      </div>
    </div>
  </template>
}

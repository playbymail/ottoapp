// Copyright (c) 2025 Michael D Henderson. All rights reserved.

import Component from '@glimmer/component';
import {service} from "@ember/service";
import { on } from '@ember/modifier';
import { action } from "@ember/object";

// Requires a TailwindCSS Plus license.

import { LinkTo } from '@ember/routing';

export default class Dashboard extends Component {
  <template>
    <!-- app/components/admin/dashboard.gjs -->
    <div class="space-y-6">
      <div>
        <h2 class="text-base font-semibold leading-7 text-gray-900">Admin Dashboard</h2>
        <p class="mt-1 text-sm leading-6 text-gray-600">
          Administrative interface for system management.
        </p>
      </div>

      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <div class="overflow-hidden rounded-lg bg-white px-4 py-5 shadow sm:p-6">
          <dt class="truncate text-sm font-medium text-gray-500">Total Users</dt>
          <dd class="mt-1 text-3xl font-semibold tracking-tight text-gray-900">--</dd>
        </div>

        <div class="overflow-hidden rounded-lg bg-white px-4 py-5 shadow sm:p-6">
          <dt class="truncate text-sm font-medium text-gray-500">Active Sessions</dt>
          <dd class="mt-1 text-3xl font-semibold tracking-tight text-gray-900">--</dd>
        </div>

        <div class="overflow-hidden rounded-lg bg-white px-4 py-5 shadow sm:p-6">
          <dt class="truncate text-sm font-medium text-gray-500">System Status</dt>
          <dd class="mt-1 text-3xl font-semibold tracking-tight text-gray-900">
            <span
              class="inline-flex items-center rounded-md bg-green-50 px-2 py-1 text-sm font-medium text-green-700 ring-1 ring-inset ring-green-600/20">
              Online
            </span>
          </dd>
        </div>
      </div>

      <div class="rounded-lg bg-blue-50 p-4">
        <div class="flex">
          <div class="shrink-0">
            <svg viewBox="0 0 20 20" fill="currentColor" aria-hidden="true" class="size-5 text-blue-400">
              <path fill-rule="evenodd"
                    d="M18 10a8 8 0 1 1-16 0 8 8 0 0 1 16 0Zm-7-4a1 1 0 1 1-2 0 1 1 0 0 1 2 0ZM9 9a.75.75 0 0 0 0 1.5h.253a.25.25 0 0 1 .244.304l-.459 2.066A1.75 1.75 0 0 0 10.747 15H11a.75.75 0 0 0 0-1.5h-.253a.25.25 0 0 1-.244-.304l.459-2.066A1.75 1.75 0 0 0 9.253 9H9Z"
                    clip-rule="evenodd" />
            </svg>
          </div>
          <div class="ml-3 flex-1 md:flex md:justify-between">
            <p class="text-sm text-blue-700">
              This is the admin dashboard. Future sprints will add user management, system settings, and audit logs.
            </p>
          </div>
        </div>
      </div>
    </div>
  </template>
}

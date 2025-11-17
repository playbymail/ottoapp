// app/templates/admin/settings.gjs

import { LinkTo } from '@ember/routing';

<template>
  <main>
    <h1 class="sr-only">Account Settings</h1>

    <header class="border-b border-gray-200">
      <!-- Secondary navigation -->
      <nav class="flex overflow-x-auto py-4">
        <ul role="list" class="flex min-w-full flex-none gap-x-6 px-4 text-sm/6 font-semibold text-gray-500 sm:px-6 lg:px-8">
          <li>
            <LinkTo @route="admin.settings.account" class="" @activeClass="text-indigo-600">Account</LinkTo>
          </li>
          <li>
            <LinkTo @route="admin.settings.notifications" class="" @activeClass="text-indigo-600">Notifications</LinkTo>
          </li>
          <li>
            <LinkTo @route="admin.settings.about" class="" @activeClass="text-indigo-600">About</LinkTo>
          </li>
        </ul>
      </nav>
    </header>

    <!-- Settings content -->
    <div class="divide-y divide-gray-200">
      {{outlet}}
    </div>
  </main>
</template>

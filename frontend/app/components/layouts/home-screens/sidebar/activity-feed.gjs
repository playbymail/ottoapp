// app/components/layouts/home-screens/sidebar/activity-feed.gjs
// You must have a Tailwind Plus License to use this component.
// https://tailwindcss.com/plus/ui-blocks/application-ui/page-examples/home-screens#sidebar

import ActivityFeedItem from './activity-feed/item';

<template>
  <header class="flex items-center justify-between border-b border-gray-200 px-4 py-4 sm:px-6 sm:py-6 lg:px-8">
    <h2 class="text-base/7 font-semibold text-gray-900">Activity feed</h2>
    <a href="#" class="text-sm/6 font-semibold text-indigo-600">View all</a>
  </header>
  <ul role="list" class="divide-y divide-gray-100">
    <ActivityFeedItem
      @image="https://images.unsplash.com/photo-1519244703995-f4e0f30006d5?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
      @actor="Michael Foster"
      @time="1h"
    >
      <p class="mt-3 truncate text-sm text-gray-500">Pushed to <span class="text-gray-700">ios-app</span> (<span
        class="font-mono text-gray-700">2d89f0c8</span> on <span class="text-gray-700">main</span>)</p>
    </ActivityFeedItem>
    <ActivityFeedItem
      @image="https://images.unsplash.com/photo-1517841905240-472988babdf9?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
      @actor="Lindsay Walton"
      @time="3h"
    >
      <p class="mt-3 truncate text-sm text-gray-500">Pushed to <span class="text-gray-700">mobile-api</span> (<span
        class="font-mono text-gray-700">249df660</span> on <span class="text-gray-700">main</span>)</p>
    </ActivityFeedItem>
    <ActivityFeedItem
      @image="https://images.unsplash.com/photo-1438761681033-6461ffad8d80?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
      @actor="Courtney Henry"
      @time="12h"
    >
      <p class="mt-3 truncate text-sm text-gray-500">Pushed to <span class="text-gray-700">ios-app</span> (<span
        class="font-mono text-gray-700">11464223</span> on <span class="text-gray-700">main</span>)</p>
    </ActivityFeedItem>
    <ActivityFeedItem
      @image="https://images.unsplash.com/photo-1517365830460-955ce3ccd263?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
      @actor="Whitney Francis"
      @time="2w"
    >
      <p class="mt-3 truncate text-sm text-gray-500">Pushed to <span class="text-gray-700">ios-app</span> (<span
        class="font-mono text-gray-700">5c1fd07f</span> on <span class="text-gray-700">main</span>)</p>
    </ActivityFeedItem>
  </ul>
</template>

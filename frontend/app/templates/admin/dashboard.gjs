// app/templates/admin/dashboard.gjs

import ActivityFeed from 'frontend/components/layouts/home-screens/sidebar/activity-feed';
import Deployments from 'frontend/components/layouts/home-screens/sidebar/deployments';

<template>
  <main class="lg:pr-96">
    <Deployments />
  </main>

  <!-- Activity feed -->
  <aside
    class="bg-gray-50 lg:fixed lg:top-16 lg:right-0 lg:bottom-0 lg:w-96 lg:overflow-y-auto lg:border-l lg:border-gray-200">
    <ActivityFeed />
  </aside>
</template>

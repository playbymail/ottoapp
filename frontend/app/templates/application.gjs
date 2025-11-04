import { service } from '@ember/service';

// https://tailwindcss.com/plus/ui-blocks/application-ui/application-shells/sidebar#sidebar-with-header
// Requires a TailwindCSS Plus license.

import Sidebar from 'frontend/components/sidebar';
import StickyHeader from 'frontend/components/sticky-header';

<template>
  <Sidebar />

  <div class="lg:pl-72">
    <StickyHeader />

    <main class="py-10">
      <div class="px-4 sm:px-6 lg:px-8">
        {{outlet}}
      </div>
    </main>
  </div>
</template>

// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Authenticated user shell with sidebar and header
// https://tailwindcss.com/plus/ui-blocks/application-ui/application-shells/sidebar#sidebar-with-header

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

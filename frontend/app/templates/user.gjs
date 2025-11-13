// Copyright (c) 2025 Michael D Henderson. All rights reserved.

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

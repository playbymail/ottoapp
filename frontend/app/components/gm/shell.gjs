// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// app/components/gm/shell.gjs

import Sidebar from './sidebar';
import StickyHeader from './sticky-header';

<template>
  <Sidebar @adminShell={{false}} @gmShell={{true}} @userShell={{false}}/>

  <div class="lg:pl-72">
    <StickyHeader />

    <main class="py-10">
      <div class="px-4 sm:px-6 lg:px-8">
        {{outlet}}
      </div>
    </main>
  </div>
</template>

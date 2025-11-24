// app/components/user/shell.gjs

import Sidebar from './sidebar';
import StickyHeader from './sticky-header';

<template>
  <Sidebar @adminShell={{false}} @gmShell={{false}} @userShell={{true}} />

  <div class="lg:pl-72">
    <StickyHeader />

    <main class="py-10">
      <div class="px-4 sm:px-6 lg:px-8">
        {{outlet}}
      </div>
    </main>
  </div>
</template>

import { service } from '@ember/service';

import ApplicationHeader from 'frontend/components/application-header';

<template>
  {{#let (service "session") as |session|}}
    <ApplicationHeader
      @session={{session}}
      @onLogout={{this.logout}}
    />
  {{/let}}

  <main class="p-6">
    {{outlet}}
  </main>
</template>

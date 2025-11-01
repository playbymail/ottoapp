import ApplicationHeader from 'frontend/components/application-header';
import { service } from '@ember/service';

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

import ApplicationHeader from '../components/application-header';
import { service } from '@ember/service';

<template>
  {{#let (service "session") (service "currentUser") as |session currentUser|}}
    <ApplicationHeader 
      @session={{session}} 
      @currentUser={{currentUser}} 
      @onLogout={{this.logout}} 
    />
  {{/let}}

  <main class="p-6">
    {{outlet}}
  </main>
</template>

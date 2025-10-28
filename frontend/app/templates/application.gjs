import { pageTitle } from 'ember-page-title';

<template>
<header class="p-4 bg-gray-100 border-b flex justify-between">
  <h1 class="font-bold text-lg">Frontend Demo</h1>

  {{#if this.session.isAuthenticated}}
    <div class="flex items-center space-x-2">
      {{#if this.currentUser.user}}
        <span>👋 {{this.currentUser.user.username}}</span>
      {{/if}}
      <button type="button" {{on "click" this.logout}}>
        Logout
      </button>
    </div>
  {{else}}
    <a href={{this.router.urlFor "login"}} class="underline text-blue-600">
      Login
    </a>
  {{/if}}
</header>

<main class="p-6">
  {{outlet}}
</main>
</template>

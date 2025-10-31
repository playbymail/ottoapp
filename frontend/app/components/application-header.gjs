import { on } from '@ember/modifier';
import { LinkTo } from '@ember/routing';

<template>
  <header class="p-4 bg-gray-100 border-b flex justify-between">
    <h1 class="font-bold text-lg">Frontend Demo</h1>

    {{#if @session.isAuthenticated}}
      <div class="flex items-center space-x-2">
        {{#if @currentUser.user}}
          <span>👋 {{@currentUser.user.username}}</span>
        {{/if}}
        <button type="button" {{on "click" @onLogout}}>
          Logout
        </button>
      </div>
    {{else}}
      <LinkTo @route="login" class="underline text-blue-600">
        Login
      </LinkTo>
    {{/if}}
  </header>
</template>

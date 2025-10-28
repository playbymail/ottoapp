import { pageTitle } from 'ember-page-title';

<template>
<h1>Sign in</h1>

{{#if this.error}}
  <p role="alert">{{this.error}}</p>
{{/if}}

<form {{on "submit" this.submit}}>
  <label>
    Username
    <input type="text" value={{this.username}} {{on "input" this.updateUsername}} />
  </label>

  <label>
    Password
    <input type="password" value={{this.password}} {{on "input" this.updatePassword}} />
  </label>

  <button type="submit">Sign in</button>
</form>
</template>

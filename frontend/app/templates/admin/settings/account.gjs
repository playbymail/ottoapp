// app/templates/admin/settings/account.gjs

import {pageTitle} from 'ember-page-title';

import SettingsAccountTab from 'frontend/components/settings/account';

<template>
  {{pageTitle "Account Settings"}}

  <SettingsAccountTab @user={{@model}} />
</template>

// app/templates/admin/settings/about.gjs

import {pageTitle} from 'ember-page-title';

import SettingsAboutTab from 'frontend/components/settings/about';

<template>
  {{pageTitle "About OttoMap"}}

  <SettingsAboutTab @version={{@model}} />
</template>

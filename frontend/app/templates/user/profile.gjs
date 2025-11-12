// app/templates/user/profile.gjs
import Profile from 'frontend/components/user/profile';
import DebugProbe from 'frontend/components/_debug/probe';
import controllerFor from "@ember/routing/lib/controller_for";

<template>
  <DebugProbe @value={{@controller}} />

  <Profile
    @model={{@model}}
    @onSave={{@controller.updateProfile}}       {{! <- bare identifier, not this.updateProfile }}
    @isSaving={{this.isSaving}}          {{! <- bare identifier, not this.isSaving }}
    @errorMessage={{this.errorMessage}}
    @successMessage={{this.successMessage}}
  />
</template>

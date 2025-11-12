// app/templates/user/profile.gjs
import Profile from 'frontend/components/user/profile';

<template>
  <Profile
    @model={{@model}}
    @onSave={{@controller.updateProfile}}
    @isSaving={{@controller.isSaving}}
    @errorMessage={{@controller.errorMessage}}
    @successMessage={{@controller.successMessage}}
  />
</template>

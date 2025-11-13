// app/templates/user/profile.gjs
import UserForm from 'frontend/components/user/form';

<template>
  <UserForm
    @model={{@model}}
    @onSave={{@controller.updateProfile}}
    @onCancel={{@controller.cancel}}
    @isSaving={{@controller.isSaving}}
    @errorMessage={{@controller.errorMessage}}
    @successMessage={{@controller.successMessage}}
  />
</template>

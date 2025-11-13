// app/components/user/profile.gjs
import UserForm from 'frontend/components/user/form';

<template>
  <UserForm
    @model={{@model}}
    @onSave={{@onSave}}
    @onCancel={{@onCancel}}
    @isSaving={{@isSaving}}
    @errorMessage={{@errorMessage}}
    @successMessage={{@successMessage}}
    @canEditUsername={{false}}
  />
</template>

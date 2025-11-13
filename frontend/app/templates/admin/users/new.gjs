// app/templates/admin/users/new.gjs
import UserForm from 'frontend/components/user/form';

<template>
  <UserForm
    @model={{@model}}
    @canEditUsername={{true}}
    @onSave={{@controller.createUser}}
    @onCancel={{@controller.cancel}}
    @isSaving={{@controller.isSaving}}
    @errorMessage={{@controller.errorMessage}}
    @successMessage={{@controller.successMessage}}
  />
</template>

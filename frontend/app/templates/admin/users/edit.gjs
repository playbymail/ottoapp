// app/templates/admin/users/edit.gjs
import UserForm from 'frontend/components/user/form';

<template>
  <UserForm
    @model={{@model}}
    @onSave={{@controller.updateUser}}
    @onCancel={{@controller.cancel}}
    @isSaving={{@controller.isSaving}}
    @errorMessage={{@controller.errorMessage}}
    @successMessage={{@controller.successMessage}}
    @canEditUsername={{true}}
  />
</template>

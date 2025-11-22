// app/controllers/admin/users/edit/profile.js

import Controller from '@ember/controller';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class AdminUsersEditController extends Controller {
  @service router;

  @tracked errorMessage = null;
  @tracked successMessage = null;
  @tracked isSaving = false;

  @action
  async updateUser() {
    this.errorMessage = null;
    this.successMessage = null;
    this.isSaving = true;

    try {
      if (!this.model.hasDirtyAttributes) {
        return;
      }

      await this.model.save();
      this.successMessage = 'User updated successfully';
    } catch (err) {
      this.model.rollbackAttributes();
      this.errorMessage = err?.message || 'An error occurred while updating the user';
    } finally {
      this.isSaving = false;
    }
  }

  @action
  cancel() {
    this.router.transitionTo('admin.users.index');
  }
}

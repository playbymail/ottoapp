import Controller from '@ember/controller';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class AdminUsersEditController extends Controller {
  @service api;
  @service router;

  @tracked errorMessage = null;
  @tracked successMessage = null;
  @tracked isSaving = false;
  @tracked showResetPasswordDialog = false;
  @tracked tempPassword = null;

  @action
  async updateUser(event) {
    event.preventDefault();
    this.errorMessage = null;
    this.successMessage = null;
    this.isSaving = true;

    const formData = new FormData(event.target);
    const userData = {
      username: formData.get('username'),
      email: formData.get('email'),
      timezone: formData.get('timezone'),
    };

    try {
      const response = await this.api.updateUser(this.model.id, userData);

      if (response.ok) {
        this.successMessage = 'User updated successfully';
        // Refresh the model
        const updatedResponse = await this.api.getUser(this.model.id);
        if (updatedResponse.ok) {
          this.model = await updatedResponse.json();
        }
      } else {
        const error = await response.json();
        this.errorMessage = error.message || 'Failed to update user';
      }
    } catch (error) {
      this.errorMessage = 'An error occurred while updating the user';
    } finally {
      this.isSaving = false;
    }
  }

  @action
  async resetPassword() {
    this.errorMessage = null;
    this.isSaving = true;

    try {
      const response = await this.api.resetPassword(this.model.id);

      if (response.ok) {
        const result = await response.json();
        this.tempPassword = result.tempPassword;
        this.showResetPasswordDialog = true;
      } else {
        const error = await response.json();
        this.errorMessage = error.message || 'Failed to reset password';
      }
    } catch (error) {
      this.errorMessage = 'An error occurred while resetting the password';
    } finally {
      this.isSaving = false;
    }
  }

  @action
  closeResetPasswordDialog() {
    this.showResetPasswordDialog = false;
    this.tempPassword = null;
  }

  @action
  cancel() {
    this.router.transitionTo('admin.users.index');
  }
}

import Controller from '@ember/controller';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class UsersPasswordController extends Controller {
  @service api;
  @service session;

  @tracked errorMessage = null;
  @tracked successMessage = null;
  @tracked isSaving = false;

  @action
  async changePassword(event) {
    event.preventDefault();
    this.errorMessage = null;
    this.successMessage = null;
    this.isSaving = true;

    const formData = new FormData(event.target);
    const currentPassword = formData.get('currentPassword');
    const newPassword = formData.get('newPassword');
    const confirmPassword = formData.get('confirmPassword');

    // Validate passwords match
    if (newPassword !== confirmPassword) {
      this.errorMessage = 'New passwords do not match';
      this.isSaving = false;
      return;
    }

    try {
      const response = await this.api.updatePassword(
        this.model.userId,
        currentPassword,
        newPassword
      );

      if (response.ok) {
        this.successMessage = 'Password changed successfully';
        event.target.reset();
      } else {
        const error = await response.json();
        this.errorMessage = error.message || 'Failed to change password';
      }
    } catch (error) {
      this.errorMessage = 'An error occurred while changing your password';
    } finally {
      this.isSaving = false;
    }
  }
}

import Controller from '@ember/controller';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class UsersProfileController extends Controller {
  @service api;
  @service router;
  @service session;

  @tracked errorMessage = null;
  @tracked successMessage = null;
  @tracked isSaving = false;

  @action
  async updateProfile(event) {
    event.preventDefault();
    this.errorMessage = null;
    this.successMessage = null;
    this.isSaving = true;

    const formData = new FormData(event.target);
    const email = formData.get('email');
    const timezone = formData.get('timezone');

    try {
      const response = await this.api.updateUser(this.model.id, {
        email,
        timezone,
      });

      if (response.ok) {
        this.successMessage = 'Profile updated successfully';
        // Refresh the model
        const updatedResponse = await this.api.getUserProfile();
        if (updatedResponse.ok) {
          this.model = await updatedResponse.json();
        }
      } else {
        const error = await response.json();
        this.errorMessage = error.message || 'Failed to update profile';
      }
    } catch (error) {
      this.errorMessage = 'An error occurred while updating your profile';
    } finally {
      this.isSaving = false;
    }
  }
}

// app/controllers/user/profile.js
import Controller from '@ember/controller';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class UserProfileController extends Controller {
  @tracked errorMessage = null;
  @tracked successMessage = null;
  @tracked isSaving = false;

  /**
   * DDAU: this action expects the component to update the model directly.
   */
  @action
  async updateProfile() {
    this.errorMessage = null;
    this.successMessage = null;
    this.isSaving = true;

    try {
      if (!this.model.hasDirtyAttributes) {
        return;
      }

      await this.model.save();
      this.successMessage = 'Profile updated successfully';
    } catch (err) {
      this.model.rollbackAttributes();
      this.errorMessage = err?.message || 'An error occurred while updating your profile';
    } finally {
      this.isSaving = false;
    }
  }
}

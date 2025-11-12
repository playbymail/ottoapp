// app/controllers/user/profile.js
import Controller from '@ember/controller';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class UserProfileController extends Controller {
  @service api;
  // keep these if you later need redirect-after-save
  // @service router;
  // @service session;

  isExpanded = false

  @tracked errorMessage = null;
  @tracked successMessage = null;
  @tracked isSaving = false;

  fiii = () => {
    console.log('app/controllers/user/profile', 'fiii called');
    return false;
  }

  /**
   * DDAU: this action expects POJO attrs from the template/component,
   * not a DOM event.
   */
  @action
  async updateProfile(changes) {
    console.log('controller:updateProfile', changes);

    this.errorMessage = null;
    this.successMessage = null;
    this.isSaving = true;

    const payload = {
      email: changes.email ?? this.model.email,
      timezone: changes.timezone ?? this.model.timezone,
    };

    try {
      const updated = await this.api.updateProfile(payload); // PUT via service
      this.successMessage = 'Profile updated successfully';
      Object.assign(this.model, updated);
    } catch (err) {
      this.errorMessage = err?.message || 'An error occurred while updating your profile';
    } finally {
      this.isSaving = false;
    }
  }

  @action
  async updateProfilex(changes) {
    console.log('app/controllers/user/profile', 'updateProfile', changes);
    this.errorMessage = null;
    this.successMessage = null;
    this.isSaving = true;

    const payload = {
      email: changes.email ?? this.model.email,
      timezone: changes.timezone ?? this.model.timezone,
    };

    try {
      // calls api.updateProfile → which should use putAsJSON
      const updated = await this.api.updateProfile(payload);

      this.successMessage = 'Profile updated successfully';

      // mutate the existing model so templates update
      Object.assign(this.model, updated);
    } catch (err) {
      this.errorMessage =
        err?.message || 'An error occurred while updating your profile';
    } finally {
      this.isSaving = false;
    }
  }

  // Arrow field -> concrete, already-bound function value
  saveProfile = async (changes) => {
    console.log('controller:saveProfile', changes);

    this.errorMessage = null;
    this.successMessage = null;
    this.isSaving = true;

    const payload = {
      email: changes.email ?? this.model.email,
      timezone: changes.timezone ?? this.model.timezone,
    };

    try {
      const updated = await this.api.updateProfile(payload); // PUT via putAsJSON
      this.successMessage = 'Profile updated successfully';
      Object.assign(this.model, updated); // mutate existing model
    } catch (err) {
      this.errorMessage = err?.message || 'An error occurred while updating your profile';
    } finally {
      this.isSaving = false;
    }
  };
}

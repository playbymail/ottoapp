import Controller from '@ember/controller';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class AdminUsersNewController extends Controller {
  @service api;
  @service router;

  @tracked errorMessage = null;
  @tracked isSaving = false;

  @action
  async createUser(event) {
    event.preventDefault();
    this.errorMessage = null;
    this.isSaving = true;

    const formData = new FormData(event.target);
    const userData = {
      username: formData.get('username'),
      email: formData.get('email'),
      password: formData.get('password'),
      timezone: formData.get('timezone'),
      roles: formData.getAll('roles'),
    };

    try {
      const response = await this.api.createUser(userData);

      if (response.ok) {
        this.router.transitionTo('admin.users.index');
      } else {
        const error = await response.json();
        this.errorMessage = error.message || 'Failed to create user';
      }
    } catch (error) {
      this.errorMessage = 'An error occurred while creating the user';
    } finally {
      this.isSaving = false;
    }
  }

  @action
  cancel() {
    this.router.transitionTo('admin.users.index');
  }
}

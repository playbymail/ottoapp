import Controller from '@ember/controller';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class AdminUsersNewController extends Controller {
  @service router;

  @tracked errorMessage = null;
  @tracked isSaving = false;

  @action
  async createUser() {
    this.errorMessage = null;
    this.isSaving = true;

    try {
      await this.model.save();
      this.router.transitionTo('admin.users.index');
    } catch (err) {
      this.errorMessage = err?.message || 'An error occurred while creating the user';
    } finally {
      this.isSaving = false;
    }
  }

  @action
  cancel() {
    this.model.rollbackAttributes();
    this.router.transitionTo('admin.users.index');
  }
}

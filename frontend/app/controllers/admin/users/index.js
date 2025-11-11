import Controller from '@ember/controller';
import { service } from '@ember/service';
import { action } from '@ember/object';

export default class AdminUsersIndexController extends Controller {
  @service router;

  @action
  editUser(userId) {
    this.router.transitionTo('admin.users.edit', userId);
  }

  @action
  createUser() {
    this.router.transitionTo('admin.users.new');
  }
}

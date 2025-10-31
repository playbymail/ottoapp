import Controller from '@ember/controller';
import { service } from '@ember/service';
import { action } from '@ember/object';

export default class ApplicationController extends Controller {
  @service session;
  @service currentUser;
  @service router;

  @action async logout() {
    await this.session.invalidate();
    this.currentUser.user = null;
    this.router.transitionTo('login');
  }
}

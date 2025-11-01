import Controller from '@ember/controller';
import { service } from '@ember/service';
import { action } from '@ember/object';

export default class ApplicationController extends Controller {
  @service session;
  @service router;

  @action async logout() {
    console.log('async logout() called');
    await this.session.invalidate();
    this.router.transitionTo('login');
  }
}

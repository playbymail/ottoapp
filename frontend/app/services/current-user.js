import Service from '@ember/service';
import { tracked } from '@glimmer/tracking';

export default class CurrentUserService extends Service {
  @tracked user = null;

  async load() {
    const res = await fetch('/api/me');
    this.user = res.ok ? await res.json() : null;
  }
}

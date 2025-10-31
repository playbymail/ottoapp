import Service from '@ember/service';
import { tracked } from '@glimmer/tracking';

export default class ApiService extends Service {
  @tracked csrf = null;

  async ensureCsrf() {
    if (!this.csrf) {
      const res = await fetch('/api/session');
      if (res.ok) this.csrf = (await res.json()).csrf;
    }
    return this.csrf;
  }

  async post(url, body) {
    const token = await this.ensureCsrf();
    return fetch(url, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', 'X-CSRF-Token': token },
      body: JSON.stringify(body),
    });
  }
}

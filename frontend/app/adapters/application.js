// app/adapters/application.js
import JSONAPIAdapter from '@ember-data/adapter/json-api';
import { service } from '@ember/service';

export default class ApplicationAdapter extends JSONAPIAdapter {
  @service session;

  namespace = 'api';

  get headers() {
    const headers = {};

    // Add CSRF token if we have a session
    if (this.session.isAuthenticated && this.session.data.authenticated.csrf) {
      headers['X-CSRF-Token'] = this.session.data.authenticated.csrf;
    }

    return headers;
  }
}

// mdhender: not used?

import Base from 'ember-simple-auth/authenticators/base';

export default class CookieAuthenticator extends Base {
  // authenticate makes a request to the /api/login route with the
  // user's name and password fields from the form data. the server
  // authenticates the credentials and returns an HTTP 200 response
  // if they are valid, anything else is taken as invalid.
  async authenticate({ username, password }) {
    console.log('esa', 'app/authenticators/cookie:restore');
    const res = await fetch('/api/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password }),
    });
    console.log('esa', 'app/authenticators/cookie:authenticate', 'res.ok', res.ok);
    if (!res.ok) throw new Error('Bad credentials');
    return { isAuthenticated: true };
  }

  async restore(data) {
    console.log('esa', 'app/authenticators/cookie:restore');
    const res = await fetch('/api/session');
    console.log('esa', 'app/authenticators/cookie:restore', 'res.ok', res.ok);
    if (res.ok) return data;
    throw new Error('Session expired');
  }

  async invalidate() {
    console.log('esa', 'app/authenticators/cookie:invalidate');
    await fetch('/api/logout', { method: 'POST' });
  }
}

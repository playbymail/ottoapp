import Base from 'ember-simple-auth/authenticators/base';

export default class CookieAuthenticator extends Base {
  async authenticate({ username, password }) {
    const res = await fetch('/api/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password }),
    });
    if (!res.ok) throw new Error('Bad credentials');
    return { isAuthenticated: true };
  }

  async restore(data) {
    const res = await fetch('/api/session');
    if (res.ok) return data;
    throw new Error('Session expired');
  }

  async invalidate() {
    await fetch('/api/logout', { method: 'POST' });
  }
}

// app/authenticators/server.js
import Base from 'ember-simple-auth/authenticators/base';

export default class ServerAuthenticator extends Base {

  async restore() {
    // console.log('esa', 'app/authenticators/server:restore');
    const resp = await fetch('/api/session', { credentials: 'include' });
    if (!resp.ok) {
      throw 'not authenticated';
    }
    const json = await resp.json();
    // ESA will persist this whole object
    // console.log('esa', 'app/authenticators/server:restore', json);
    return json;
  }

  async authenticate(credentials) {
    console.log('esa', 'app/authenticators/server:authenticate');
    // hit your POST /api/login which sets the cookie
    const resp = await fetch('/api/login', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(credentials),
    });
    if (!resp.ok) {
      throw 'bad credentials';
    }

    // after login, load the session details
    return this.restore();
  }

  async invalidate() {
    // console.log('esa', 'app/authenticators/server:invalidate', 'enter');
    try {
      await fetch('/api/logout', {
        method: 'POST',
        credentials: 'include',
      });
      // even if the request never comes back, we still want to clear the Ember session
      // console.log('esa', 'app/authenticators/server:invalidate', 'fetch finished');
    } catch {
      // ignore network/abort errors on logout.
    }
    // console.log('esa', 'app/authenticators/server:invalidate', 'exit');
  }
}


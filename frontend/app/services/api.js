// app/services/api.js
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

  // small helper to pull message from JSON error bodies
  async _errorFromResponse(res, fallback) {
    try {
      const data = await res.json();
      return new Error(data.message || fallback);
    } catch (_) {
      return new Error(fallback);
    }
  }

  // implement HTTP request method helpers
  async get(url) {
    return fetch(url, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
    });
  }

  // NEW: get but always return JSON or throw
  async getAsJSON(url) {
    const res = await fetch(url, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'same-origin', // match what your route had
    });
    if (!res.ok) {
      // you can make a better error object here
      throw new Error(`GET ${url} failed with ${res.status}`);
    }
    return res.json();
  }

  async post(url, body) {
    const token = await this.ensureCsrf();
    return fetch(url, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', 'X-CSRF-Token': token },
      body: JSON.stringify(body),
    });
  }

  // NEW: post but always return JSON or throw
  async postAsJSON(url, body) {
    const token = await this.ensureCsrf();
    const res = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-CSRF-Token': token,
      },
      credentials: 'same-origin',
      body: JSON.stringify(body),
    });
    if (!res.ok) {
      // you can make a better error object here
      throw await this._errorFromResponse(res, `POST ${url} failed`);
    }
    return res.json();
  }

  async patch(url, body) {
    const token = await this.ensureCsrf();
    return fetch(url, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json', 'X-CSRF-Token': token },
      body: JSON.stringify(body),
    });
  }

  async put(url, body) {
    const token = await this.ensureCsrf();
    return fetch(url, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', 'X-CSRF-Token': token },
      body: JSON.stringify(body),
    });
  }

  // NEW: put but always return JSON or throw — this is what the controller wants
  async putAsJSON(url, body) {
    const token = await this.ensureCsrf();
    const res = await fetch(url, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'X-CSRF-Token': token,
      },
      credentials: 'same-origin',
      body: JSON.stringify(body),
    });
    if (!res.ok) {
      // you can make a better error object here
      throw await this._errorFromResponse(res, `PUT ${url} failed`);
    }
    return res.json();
  }

  // User Management API Methods

  async getUserProfile() {
    return this.get('/api/users/me');
  }

  async getUsers() {
    return this.get('/api/users');
  }

  async getUser(userId) {
    return this.get(`/api/users/${userId}`);
  }

  async createUser(userData) {
    return this.post('/api/users', userData);
  }

  async updateUser(userId, userData) {
    return this.patch(`/api/users/${userId}`, userData);
  }

  async updatePassword(userId, currentPassword, newPassword) {
    return this.put(`/api/users/${userId}/password`, {
      currentPassword,
      newPassword,
    });
  }

  async resetPassword(userId) {
    return this.post(`/api/users/${userId}/reset-password`, {});
  }

  async updateUserRole(userId, rolesToAdd, rolesToRemove) {
    return this.patch(`/api/users/${userId}/role`, {
      add: rolesToAdd,
      remove: rolesToRemove,
    });
  }

  async getProfile() {
    console.log('app/services/api', 'getProfile');
    return this.getAsJSON('/api/profile');
  }

  async updateProfile(data) {
    console.log('app/services/api', 'change updateProfile to PUT not POST');
    return this.postAsJSON('/api/profile', data);
  }

  async getTimezones() {
    return this.get('/api/timezones');
  }
}

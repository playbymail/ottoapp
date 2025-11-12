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

  async get(url) {
    return fetch(url, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
    });
  }

  async post(url, body) {
    const token = await this.ensureCsrf();
    return fetch(url, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', 'X-CSRF-Token': token },
      body: JSON.stringify(body),
    });
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
    return this.get('/api/profile');
  }

  async updateProfile(data) {
    return this.post('/api/profile', data);
  }

  async getTimezones() {
    return this.get('/api/timezones');
  }
}

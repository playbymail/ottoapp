import Service from 'ember-simple-auth/services/session';

export default class SessionService extends Service {
  async handleInvalidation(routeAfterInvalidation) {
    // i can't figure out how to set routeAfterInvalidation with ESA config,
    // so i'm whacking it here.
    routeAfterInvalidation = '/login'

    // Perform any custom cleanup or server-side invalidation here

    // Call the super method to ensure default Ember Simple Auth invalidation logic runs
    await super.handleInvalidation(routeAfterInvalidation);
  }

  /**
   * Get the current user's roles
   * @returns {Array<string>} Array of role names
   */
  get roles() {
    return this.data?.authenticated?.user?.roles || [];
  }

  /**
   * Check if the current user has a specific role
   * @param {string} role - Role name to check
   * @returns {boolean}
   */
  hasRole(role) {
    return this.roles.includes(role);
  }

  /**
   * Check if the current user is an admin
   * @returns {boolean}
   */
  get isAdmin() {
    return this.hasRole('admin');
  }

  /**
   * Check if the current user is a regular user
   * @returns {boolean}
   */
  get isUser() {
    return this.hasRole('user');
  }

  /**
   * Check if the current user can access admin routes
   * @returns {boolean}
   */
  get canAccessAdminRoutes() {
    return this.isAdmin;
  }

  /**
   * Check if the current user can access user routes
   * @returns {boolean}
   */
  get canAccessUserRoutes() {
    return this.isUser || this.isAdmin;
  }

  /**
   * Get current user ID
   * @returns {number|null}
   */
  get currentUserId() {
    return this.data?.authenticated?.user?.id || null;
  }

  /**
   * Get current user's permissions
   * @returns {Object}
   */
  get permissions() {
    return this.data?.authenticated?.user?.permissions || {};
  }

  /**
   * Check if current user can edit usernames
   * @returns {boolean}
   */
  get canEditUsername() {
    return this.permissions.canEditUsername || false;
  }
}

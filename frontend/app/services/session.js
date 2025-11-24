// app/services/session.js

import Service from 'ember-simple-auth/services/session';

export default class SessionService extends Service {
  async handleInvalidation(routeAfterInvalidation) {
    // I can't figure out how to set routeAfterInvalidation with ESA config,
    // so I'm whacking it here.
    routeAfterInvalidation = '/login'

    // Perform any custom cleanup or server-side invalidation here

    // Call the super method to ensure default Ember Simple Auth invalidation logic runs
    await super.handleInvalidation(routeAfterInvalidation);
  }

  /**
   * Check if the current user can access admin routes
   * @returns {boolean}
   */
  get canAccessAdminRoutes() {
    return this.data?.authenticated?.roles.accessAdminRoutes;
  }

  /**
   * Check if the current user can access gm routes
   * @returns {boolean}
   */
  get canAccessGMRoutes() {
    return this.data?.authenticated?.roles.accessGMRoutes;
  }

  /**
   * Check if the current user can access user routes
   * @returns {boolean}
   */
  get canAccessUserRoutes() {
    return this.data?.authenticated?.roles.accessUserRoutes;
  }

  /**
   * Check if current user can edit handles
   * @returns {boolean}
   */
  get canEditHandle() {
    return this.data?.authenticated?.roles.editHandle;
  }

  /**
   * Get current user handle
   * @returns {string}
   */
  get getHandle() {
    return this.data?.authenticated?.handle || "guest";
  }

  /**
   * Get current user ID
   * @returns {number|null}
   */
  get getUserID() {
    return this.data?.authenticated?.userId || 0;
  }
}

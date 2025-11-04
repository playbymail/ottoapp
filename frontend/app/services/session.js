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
}

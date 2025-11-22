// app/routes/user/maps.js
import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class UserMapsRoute extends Route {
  @service store;

  async model() {
    const documents = await this.store.query('document', {
      filter: {
        kind: 'worldographer-map',
      },
    });

    // Sort by updatedAt descending (newest first)
    return documents.slice().sort((a, b) => {
      return new Date(b.updatedAt) - new Date(a.updatedAt);
    });
  }
}

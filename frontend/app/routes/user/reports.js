// app/routes/user/reports.js
import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class UserReportsRoute extends Route {
  @service store;

  async model() {
    const documents = await this.store.query('document', {
      filter: {
        kind: 'turn-report-file',
      },
    });

    // Sort by updatedAt descending (newest first)
    return documents.slice().sort((a, b) => {
      return new Date(b.updatedAt) - new Date(a.updatedAt);
    });
  }
}

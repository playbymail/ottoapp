// app/routes/user/extracts/show.js

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class UserExtractsShowRoute extends Route {
  @service store;
  @service router;

  async model(params) {
    try {
      let document = await this.store.findRecord('document', params.document_id);

      // Verify type
      if (!document.isTurnReportExtract) {
        this.router.transitionTo('user.extracts');
        return;
      }

      // Fetch content
      let response = await fetch(document.downloadUrl);
      if (!response.ok) {
        throw new Error(`Failed to fetch content: ${response.statusText}`);
      }
      let text = await response.text();

      return {
        document,
        content: text,
      };
    } catch (e) {
      console.log('app/routes/user/extracts/show.js', 'error loading extract', e);
      this.router.transitionTo('user.extracts');
    }
  }
}

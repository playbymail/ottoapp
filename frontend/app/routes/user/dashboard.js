// app/routes/user/dashboard.js

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class UserDashboardRoute extends Route {
  @service session;
  @service store;

  async model() {
    // You can tweak the query params as needed (filter by clan, game, etc.)
    let documents = await this.store.query('document', {});

    // Use .slice() to convert the RecordArray to a native JS array
    let docs = documents.slice();
    // console.log('app/routes/user/dashboard', 'docs', docs);

    // Recent map files
    let recentMapFiles = docs
      .filter((doc) => doc.isMap)
      .sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt))
      .slice(0, 5);

    // Recent turn report files
    let recentTurnReportFiles = docs
      .filter((doc) => doc.isTurnReportFile)
      .sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt))
      .slice(0, 5);
    // console.log('app/routes/user/dashboard', 'recentTurnReportFiles', recentTurnReportFiles);

    // Recent turn report extracts
    let recentTurnReportExtracts = docs
      .filter((doc) => doc.isTurnReportExtract)
      .sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt))
      .slice(0, 5);
    // console.log('app/routes/user/dashboard', 'recentTurnReportExtracts', recentTurnReportExtracts);

    return {
      recentMapFiles,
      recentTurnReportFiles,
      recentTurnReportExtracts,
    };
  }
}

// app/routes/gm/turn-report-files/upload.js

import Route from '@ember/routing/route';

export default class UserDocumentsUploadRoute extends Route {
  async model() {
    // if you know the current clan/game, you can include it here
    console.log('app/routes/gm/turn-report-files/upload');
    return {
      gameId: '0301',
    };
  }
}

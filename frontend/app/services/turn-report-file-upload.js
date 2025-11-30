// app/services/turn-report-files-upload.js

import Service from '@ember/service';
import { service } from '@ember/service';

export default class TurnReportFileUploadService extends Service {
  @service session;

  constructor(...args) {
    console.log('app/services/turn-report-files-upload');
    super(...args);
  }

  async uploadTurnReportFile(gameCode, file) {
    console.log('app/services/turn-report-files-upload', 'uploadTurnReportFile', gameCode);
    let formData = new FormData();
    formData.append('file', file);

    const apiUrl = `/api/games/${encodeURIComponent(gameCode)}/turn-report-files`;
    console.log('app/services/turn-report-files-upload', 'uploadTurnReportFile', apiUrl);
    const response = await fetch(apiUrl, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: formData,
    });

    if (!response.ok) {
      try {
        let json = await response.json();
        let first = json?.errors?.[0];
        throw new Error(first?.detail || first?.title || 'Upload failed.');
      } catch {
        throw new Error('Server error during upload.');
      }
    }

    let json = await response.json();
    let doc = json.data;
    let documentName = doc?.attributes?.['document-name'];

    return { documentName, raw: json };
  }
}

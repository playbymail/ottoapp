// app/controllers/gm/turn-report-files/upload.js

import Controller from '@ember/controller';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { service } from '@ember/service';

export default class GmTurnReportFilesUploadController extends Controller {
  @service turnReportFileUpload; // custom service

  @tracked isUploading = false;
  @tracked errorMessage = null;
  @tracked successMessage = null;

  constructor(...args) {
    console.log('app/controllers/gm/turn-report-files/upload');
    super(...args);
  }

  @action
  async handleFileSelected(file) {
    console.log('app/controllers/gm/turn-report-files/upload', 'handleFileSelected');
    this.errorMessage = null;
    this.successMessage = null;

    if (!file) {
      console.log('app/controllers/gm/turn-report-files/upload', 'handleFileSelected', '!file');
      return;
    }

    this.isUploading = true;
    try {
      console.log('app/controllers/gm/turn-report-files/upload', 'handleFileSelected', 'try', this.turnReportFileUpload);
      console.log('app/controllers/gm/turn-report-files/upload', 'handleFileSelected', 'try', this.turnReportFileUpload.uploadTurnReportFile);
      // This will:
      // - enforce gm role on backend (401/403 on failure)
      // - enforce Word mime/size
      // - parse + validate clan heading
      // - save to documents table with derived name
      let result = await this.turnReportFileUpload.uploadTurnReportFile('0301', file);

      // You can tailor this to whatever JSON:API response you send back
      this.successMessage = `Uploaded ${file.name} as ${result.documentName}`;
    } catch (e) {
      // Map network / validation / JSON:API errors to a friendly message
      this.errorMessage = e?.message ?? 'Upload failed. Please try again.';
    } finally {
      this.isUploading = false;
    }
  }
}

// app/models/document.js

import Model, {attr} from '@ember-data/model';

export default class DocumentModel extends Model {
  @attr('string') ownerHandle;
  @attr('string') userHandle;
  @attr('string') gameId; // Game identifier (e.g., "TN3.1")
  @attr('string') clanNo; // Clan number as string ("0134")
  @attr('string') documentName; // File name with extension
  @attr('string') documentType; // "turn-report-file" | "turn-report-extract" | "worldographer-map"
  @attr('string') processingStatus; // "uploaded" | "processing" | "success" | "failed"
  @attr('boolean') isShared;
  // permissions
  @attr('boolean') canRead;
  @attr('boolean') canWrite;
  @attr('boolean') canDelete;
  @attr('boolean') canShare;

  // Timestamps
  @attr('date') createdAt;
  @attr('date') updatedAt;

  // JSON:API "links" object is passed through untouched.
  // Each is accessed like: this.links.contents.href
  @attr() links;

  //
  // Computed helpers
  //

  get isMap() {
    return this.documentType === 'worldographer-map';
  }

  get isTurnReportFile() {
    return this.documentType === 'turn-report-file';
  }

  get isTurnReportExtract() {
    return this.documentType?.startsWith('turn-report-extract');
  }

  // Success vs error boolean helpers
  get isSuccess() {
    return this.processingStatus === 'success';
  }

  get isError() {
    return this.processingStatus === 'failed' || this.processingStatus === 'error';
  }

  get isProcessing() {
    return this.processingStatus === 'processing';
  }

  get isUploaded() {
    return this.processingStatus === 'uploaded';
  }

  //
  // Convenience accessors for links
  //

  get downloadUrl() {
    return this.links?.contents?.href ?? null;
  }

  get logUrl() {
    return this.links?.log?.href ?? null;
  }

  get inputDocumentUrl() {
    return this.links?.input?.href ?? null;
  }

  get extractUrl() {
    return this.links?.extract?.href ?? null;
  }

  get outputUrl() {
    return this.links?.output?.href ?? null;
  }
}

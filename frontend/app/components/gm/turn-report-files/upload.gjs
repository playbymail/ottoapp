// app/components/gm/turn-report-files/upload.gjs

import Component from '@glimmer/component';
import { action } from '@ember/object';
import {on} from '@ember/modifier';


export default class TurnReportUpload extends Component {
  @action
  handleFileChange(event) {
    console.log('app/components/gm/turn-report-files/upload', 'handleFileChange', event);
    let [file] = event.target.files;
    if (!file) {
      console.log('app/components/gm/turn-report-files/upload', 'handleFileChange', '!file');
      return;
    }
    console.log('app/components/gm/turn-report-files/upload', 'handleFileChange', file);

    // Lightweight client-side checks (you’ll still enforce server-side)
    // Accept: Word docs, max ~150kb
    let maxSize = 150 * 1024;

    if (file.size > maxSize) {
      // If you want to handle this in the component, you can expose
      // another callback like @onValidationError.
      window.alert('File is too large. Max size is 150 KB.');
      event.target.value = '';
      return;
    }

    // Let the controller/service deal with the rest
    console.log('app/components/gm/turn-report-files/upload', 'handleFileChange', 'onFileSelected');
    this.args.onFileSelected?.(file);
  }

  get isUploading() {
    console.log('app/components/gm/turn-report-files/upload', 'handleFileChange', event);
    return this.args.isUploading;
  }

  <template>
    <div class="rounded-lg border border-dashed border-gray-300 px-6 py-10">
      <div class="text-center">
        <div class="mt-4 flex text-sm text-gray-600 justify-center">
          <label
            class="relative cursor-pointer rounded-md bg-white font-medium text-indigo-600 focus-within:outline-none focus-within:ring-2 focus-within:ring-indigo-500 focus-within:ring-offset-2 hover:text-indigo-500"
          >
            <span>Select a Word document</span>
            <input
              type="file"
              class="sr-only {{if this.isUploading "disabled"}}"
              accept=".doc,.docx,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document"
              {{on "change" this.handleFileChange}}
            />
          </label>
          <p class="pl-1">or drag and drop</p>
        </div>

        <p class="text-xs text-gray-500 mt-2">
          DOC/DOCX up to 150KB
        </p>

        {{#if @isUploading}}
          <p class="mt-4 text-sm text-gray-500">
            Uploading…
          </p>
        {{/if}}
      </div>
    </div>
  </template>
}

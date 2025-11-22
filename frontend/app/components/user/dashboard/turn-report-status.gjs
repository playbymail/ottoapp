// app/components/user/dashboard/turn-report-status.gjs

import Component from "@glimmer/component";
import { concat } from '@ember/helper';
import { LinkTo } from "@ember/routing";

import eq from 'frontend/helpers/eq';
import or from 'frontend/helpers/or';

export default class TurnReportStatus extends Component {
  <template>
    <div class="flex flex-col items-end gap-1 text-xs">
      {{!-- Status pill --}}
      <span
        class={{concat
        "inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium "
        (if (eq @status "success")
          "bg-green-100 text-green-800 "
          (if (or (eq @status "failed") (eq @status "error"))
            "bg-red-100 text-red-800 "
            (if (eq @status "processing")
              "bg-indigo-100 text-indigo-800 "
              "bg-gray-100 text-gray-800 "
            )
          )
        )
      }}
      >
        {{#if (eq @status "success")}}
          Success
        {{else if (or (eq @status "failed") (eq @status "error"))}}
          Error
        {{else if (eq @status "processing")}}
          Processing
        {{else}}
          Uploaded
        {{/if}}
      </span>

      {{!-- Stepper: Uploaded -> Processing -> (Success | Error) --}}
      <ol class="flex items-center gap-1 text-[10px] text-gray-500">
        {{!-- Uploaded --}}
        <li class="flex items-center gap-1">
          <span
            class={{concat
            "h-1.5 w-1.5 rounded-full "
            "bg-indigo-600 "
          }}
          ></span>
          <span>Uploaded</span>
        </li>

        <span class="mx-1 text-gray-400">&rarr;</span>

        {{!-- Processing --}}
        <li class="flex items-center gap-1">
          <span
            class={{concat
            "h-1.5 w-1.5 rounded-full "
            (if (or (eq @status "processing") (eq @status "success") (eq @status "failed") (eq @status "error"))
              "bg-indigo-600 "
              "bg-gray-300 "
            )
          }}
          ></span>
          <span>Processing</span>
        </li>

        <span class="mx-1 text-gray-400">&rarr;</span>

        {{!-- Final: Success or Error --}}
        <li class="flex items-center gap-1">
          <span
            class={{concat
            "h-1.5 w-1.5 rounded-full "
            (if (eq @status "success")
              "bg-green-600 "
              (if (or (eq @status "failed") (eq @status "error"))
                "bg-red-600 "
                "bg-gray-300 "
              )
            )
          }}
          ></span>
          <span>
            {{#if (or (eq @status "failed") (eq @status "error"))}}
              Error
            {{else}}
              Success
            {{/if}}
          </span>
        </li>
      </ol>
    </div>
  </template>
}

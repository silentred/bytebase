<template>
  <BBTable
    :columnList="columnList"
    :dataSource="tableList"
    :showHeader="true"
    :leftBordered="true"
    :rightBordered="true"
    :rowClickable="false"
  >
    <template v-slot:body="{ rowData: table }">
      <BBTableCell :leftPadding="4" class="w-16">
        {{ table.name }}
      </BBTableCell>
      <BBTableCell v-if="mode == 'TABLE'" class="w-8">
        {{ table.engine }}
      </BBTableCell>
      <BBTableCell v-if="mode == 'TABLE'" class="w-8">
        {{ table.collation }}
      </BBTableCell>
      <BBTableCell v-if="mode == 'TABLE'" class="w-8">
        {{ table.rowCount }}
      </BBTableCell>
      <BBTableCell v-if="mode == 'TABLE'" class="w-8">
        {{ bytesToString(table.dataSize) }}
      </BBTableCell>
      <BBTableCell v-if="mode == 'TABLE'" class="w-8">
        {{ bytesToString(table.indexSize) }}
      </BBTableCell>
      <BBTableCell v-if="mode == 'TABLE'" class="w-8">
        {{ table.dataFree }}
      </BBTableCell>
      <BBTableCell class="w-8">
        {{ table.syncStatus }}
      </BBTableCell>
      <BBTableCell class="w-16">
        {{ humanizeTs(table.lastSuccessfulSyncTs) }}
      </BBTableCell>
    </template>
  </BBTable>
</template>

<script lang="ts">
import { computed, PropType } from "vue";
import { BBTableColumn } from "../bbkit/types";
import { Table } from "../types";
import { bytesToString } from "../utils";

type Mode = "TABLE" | "VIEW";

const columnListMap: Map<Mode, BBTableColumn[]> = new Map([
  [
    "TABLE",
    [
      {
        title: "Name",
      },
      {
        title: "Engine",
      },
      {
        title: "Collation",
      },
      {
        title: "Row count est.",
      },
      {
        title: "Data size",
      },
      {
        title: "Index size",
      },
      {
        title: "Free size",
      },
      {
        title: "Sync status",
      },
      {
        title: "Last successful sync",
      },
    ],
  ],
  [
    "VIEW",
    [
      {
        title: "Name",
      },
      {
        title: "Sync status",
      },
      {
        title: "Last successful sync",
      },
    ],
  ],
]);

export default {
  name: "TableTable",
  components: {},
  props: {
    mode: {
      default: "TABLE",
      type: String as PropType<Mode>,
    },
    tableList: {
      required: true,
      type: Object as PropType<Table[]>,
    },
  },
  setup(props, ctx) {
    const columnList = computed(() => {
      return columnListMap.get(props.mode);
    });

    return {
      columnList,
      bytesToString,
    };
  },
};
</script>

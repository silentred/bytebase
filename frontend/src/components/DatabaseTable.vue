<template>
  <BBTable
    :columnList="columnList"
    :dataSource="databaseList"
    :showHeader="true"
    :leftBordered="bordered"
    :rightBordered="bordered"
    @click-row="clickDatabase"
  >
    <template v-slot:body="{ rowData: database }">
      <BBTableCell :leftPadding="4" class="w-16">
        {{ database.name }}
      </BBTableCell>
      <BBTableCell v-if="showProjectColumn" class="w-16">
        <div class="flex flex-row space-x-2 items-center">
          <div>{{ projectName(database.project) }}</div>
          <div class="tooltip-wrapper">
            <svg
              v-if="database.project.workflowType == 'UI'"
              class="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            ></svg>
            <template v-else-if="database.project.workflowType == 'VCS'">
              <span class="tooltip">Version control enabled</span>
              <svg
                class="w-4 h-4 text-control hover:text-control-hover"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
                ></path>
              </svg>
            </template>
          </div>
        </div>
      </BBTableCell>
      <BBTableCell v-if="showEnvironmentColumn" class="w-12">
        {{ environmentName(database.instance.environment) }}
      </BBTableCell>
      <BBTableCell v-if="showInstanceColumn" class="w-24">
        {{ instanceName(database.instance) }}
      </BBTableCell>
      <BBTableCell v-if="showMiscColumn" class="w-8">
        {{ database.characterSet }}
      </BBTableCell>
      <BBTableCell v-if="showMiscColumn" class="w-8">
        {{ database.collation }}
      </BBTableCell>
      <BBTableCell class="w-8">
        {{ database.syncStatus }}
      </BBTableCell>
      <BBTableCell v-if="showMiscColumn" class="w-16">
        {{ humanizeTs(database.lastSuccessfulSyncTs) }}
      </BBTableCell>
    </template>
  </BBTable>
</template>

<script lang="ts">
import { computed, PropType } from "vue";
import { useStore } from "vuex";
import { useRouter } from "vue-router";
import { databaseSlug, isDBAOrOwner } from "../utils";
import { Database } from "../types";
import { BBTableColumn } from "../bbkit/types";

type Mode = "ALL" | "ALL_SHORT" | "INSTANCE" | "PROJECT" | "PROJECT_SHORT";

const columnListMap: Map<
  Mode | "ALL_HIDE_INSTANCE" | "ALL_HIDE_INSTANCE_SHORT",
  BBTableColumn[]
> = new Map([
  [
    "ALL",
    [
      {
        title: "Name",
      },
      {
        title: "Project",
      },
      {
        title: "Environment",
      },
      {
        title: "Instance",
      },
      {
        title: "Character set",
      },
      {
        title: "Collation",
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
    "ALL_SHORT",
    [
      {
        title: "Name",
      },
      {
        title: "Project",
      },
      {
        title: "Environment",
      },
      {
        title: "Instance",
      },
      {
        title: "Sync status",
      },
    ],
  ],
  [
    "ALL_HIDE_INSTANCE",
    [
      {
        title: "Name",
      },
      {
        title: "Project",
      },
      {
        title: "Environment",
      },
      {
        title: "Character set",
      },
      {
        title: "Collation",
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
    "ALL_HIDE_INSTANCE_SHORT",
    [
      {
        title: "Name",
      },
      {
        title: "Project",
      },
      {
        title: "Environment",
      },
    ],
  ],
  [
    "INSTANCE",
    [
      {
        title: "Name",
      },
      {
        title: "Project",
      },
      {
        title: "Character set",
      },
      {
        title: "Collation",
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
    "PROJECT",
    [
      {
        title: "Name",
      },
      {
        title: "Environment",
      },
      {
        title: "Character set",
      },
      {
        title: "Collation",
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
    "PROJECT_SHORT",
    [
      {
        title: "Name",
      },
      {
        title: "Environment",
      },
      {
        title: "Sync status",
      },
    ],
  ],
]);

export default {
  name: "DatabaseTable",
  emits: ["select-database-id"],
  components: {},
  props: {
    bordered: {
      default: true,
      type: Boolean,
    },
    mode: {
      default: "ALL",
      type: String as PropType<Mode>,
    },
    singleInstance: {
      default: true,
      type: Boolean,
    },
    customClick: {
      default: false,
      type: Boolean,
    },
    databaseList: {
      required: true,
      type: Object as PropType<Database[]>,
    },
  },
  setup(props, { emit, attrs }) {
    const store = useStore();
    const router = useRouter();

    const currentUser = computed(() => store.getters["auth/currentUser"]());

    const showProjectColumn = computed(() => {
      return props.mode != "PROJECT" && props.mode != "PROJECT_SHORT";
    });

    const showEnvironmentColumn = computed(() => {
      return props.mode != "INSTANCE";
    });

    const showInstanceColumn = computed(() => {
      return (
        (props.mode == "ALL" || props.mode == "ALL_SHORT") &&
        isDBAOrOwner(currentUser.value.role)
      );
    });

    const showMiscColumn = computed(() => {
      return props.mode != "ALL_SHORT" && props.mode != "PROJECT_SHORT";
    });

    const columnList = computed(() => {
      if (
        (props.mode == "ALL" || props.mode == "ALL_SHORT") &&
        !showInstanceColumn.value
      ) {
        return props.mode == "ALL"
          ? columnListMap.get("ALL_HIDE_INSTANCE")
          : columnListMap.get("ALL_HIDE_INSTANCE_SHORT");
      }
      return columnListMap.get(props.mode);
    });

    const clickDatabase = function (section: number, row: number) {
      const database = props.databaseList[row];
      if (props.customClick) {
        emit("select-database-id", database.id);
      } else {
        router.push(`/db/${databaseSlug(database)}`);
      }
    };

    return {
      showProjectColumn,
      showEnvironmentColumn,
      showInstanceColumn,
      showMiscColumn,
      columnList,
      clickDatabase,
    };
  },
};
</script>

<template>
  <div class="flex flex-col">
    <div class="px-2 py-2 flex justify-between items-center">
      <EnvironmentTabFilter
        :selectedId="state.selectedEnvironment?.id"
        @select-environment="selectEnvironment"
      />
      <BBTableSearch
        ref="searchField"
        :placeholder="'Search issue name'"
        @change-text="(text) => changeSearchText(text)"
      />
    </div>
    <IssueTable
      :leftBordered="false"
      :rightBordered="false"
      :topBordered="true"
      :bottomBordered="true"
      :issueSectionList="[
        {
          title: 'Assigned',
          list: filteredList(state.assignedList).sort(openIssueSorter),
        },
        {
          title: 'Created',
          list: filteredList(state.createdList).sort(openIssueSorter),
        },
        {
          title: 'Subscribed',
          list: filteredList(state.subscribeList).sort(openIssueSorter),
        },
        {
          title: 'Recently Closed',
          list: filteredList(state.closeList).sort((a, b) => {
            return b.updatedTs - a.updatedTs;
          }),
        },
      ]"
    />
  </div>
</template>

<script lang="ts">
import { watchEffect, computed, onMounted, reactive, ref } from "vue";
import { useStore } from "vuex";
import { useRouter } from "vue-router";
import EnvironmentTabFilter from "../components/EnvironmentTabFilter.vue";
import IssueTable from "../components/IssueTable.vue";
import { activeEnvironment, activeTask } from "../utils";
import { Environment, Issue, TaskStatus, UNKNOWN_ID } from "../types";

interface LocalState {
  createdList: Issue[];
  assignedList: Issue[];
  subscribeList: Issue[];
  closeList: Issue[];
  searchText: string;
  selectedEnvironment?: Environment;
}

export default {
  name: "Home",
  components: {
    EnvironmentTabFilter,
    IssueTable,
  },
  props: {},
  setup(props, ctx) {
    const searchField = ref();

    const store = useStore();
    const router = useRouter();

    const state = reactive<LocalState>({
      createdList: [],
      assignedList: [],
      subscribeList: [],
      closeList: [],
      searchText: "",
      selectedEnvironment: router.currentRoute.value.query.environment
        ? store.getters["environment/environmentById"](
            router.currentRoute.value.query.environment
          )
        : undefined,
    });

    const currentUser = computed(() => store.getters["auth/currentUser"]());

    onMounted(() => {
      // Focus on the internal search field when mounted
      searchField.value.$el.querySelector("#search").focus();
    });

    const prepareIssueList = () => {
      // It will also be called when user logout
      if (currentUser.value.id != UNKNOWN_ID) {
        store
          .dispatch("issue/fetchIssueListForUser", currentUser.value.id)
          .then((issueList: Issue[]) => {
            state.assignedList = [];
            state.createdList = [];
            state.subscribeList = [];
            state.closeList = [];
            for (const issue of issueList) {
              // "OPEN"
              if (issue.status === "OPEN") {
                if (issue.assignee?.id === currentUser.value.id) {
                  state.assignedList.push(issue);
                } else if (issue.creator.id === currentUser.value.id) {
                  state.createdList.push(issue);
                } else if (
                  issue.subscriberIdList.includes(currentUser.value.id)
                ) {
                  state.subscribeList.push(issue);
                }
              }
              // "DONE" or "CANCELED"
              else if (issue.status === "DONE" || issue.status === "CANCELED") {
                if (
                  issue.creator.id === currentUser.value.id ||
                  issue.assignee?.id === currentUser.value.id ||
                  issue.subscriberIdList.includes(currentUser.value.id)
                ) {
                  state.closeList.push(issue);
                }
              }
            }
          });
      }
    };

    watchEffect(prepareIssueList);

    const selectEnvironment = (environment: Environment) => {
      state.selectedEnvironment = environment;
      if (environment) {
        router.replace({
          name: "workspace.home",
          query: { environment: environment.id },
        });
      } else {
        router.replace({ name: "workspace.home" });
      }
    };

    const changeSearchText = (searchText: string) => {
      state.searchText = searchText;
    };

    const filteredList = (list: Issue[]) => {
      if (!state.selectedEnvironment && !state.searchText) {
        // Select "All"
        return list;
      }
      return list.filter((issue) => {
        return (
          (!state.selectedEnvironment ||
            activeEnvironment(issue.pipeline).id ===
              state.selectedEnvironment.id) &&
          (!state.searchText ||
            issue.name.toLowerCase().includes(state.searchText.toLowerCase()))
        );
      });
    };

    const openIssueSorter = (a: Issue, b: Issue) => {
      const statusOrder = (status: TaskStatus) => {
        switch (status) {
          case "FAILED":
            return 0;
          case "PENDING_APPROVAL":
            return 1;
          case "PENDING":
            return 2;
          case "RUNNING":
            return 3;
          case "DONE":
            return 4;
          case "CANCELED":
            return 5;
        }
      };
      const aStatusOrder = statusOrder(activeTask(a.pipeline).status);
      const bStatusOrder = statusOrder(activeTask(b.pipeline).status);
      if (aStatusOrder == bStatusOrder) {
        return b.updatedTs - a.updatedTs;
      }
      return aStatusOrder - bStatusOrder;
    };

    return {
      searchField,
      state,
      filteredList,
      selectEnvironment,
      changeSearchText,
      openIssueSorter,
    };
  },
};
</script>

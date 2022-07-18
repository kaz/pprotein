<template>
  <section>
    <GroupEntriesTable
      :group-id="groupId"
      :entries="$store.getters.entriesByGroup(groupId)"
    />
    <AddMemo :group-id="groupId" />
  </section>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import GroupEntriesTable from "./GroupEntriesTable.vue";
import AddMemo from "./AddMemo.vue";

export default defineComponent({
  components: {
    AddMemo,
    GroupEntriesTable,
  },
  computed: {
    groupId() {
      const gid = this.$route.params.gid;
      return typeof gid == "string" ? gid : gid[0];
    },
  },
});
</script>

<style scoped lang="scss">
section {
  margin: 2em;
}
</style>

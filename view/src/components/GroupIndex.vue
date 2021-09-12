<template>
  <section>
    <GroupEntriesTable
      :group-id="$route.params.gid"
      :entries="$store.getters.entriesByGroup($route.params.gid)"
    />
  </section>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import GroupEntriesTable from "./GroupEntriesTable.vue";

export default defineComponent({
  components: {
    GroupEntriesTable,
  },
});
</script>

<style scoped lang="scss">
section {
  margin: 2em;
}
</style>

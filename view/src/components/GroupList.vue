<template>
  <section>
    <div class="control">
      <button @click="collect">Collect</button>
    </div>
    <details v-for="(group, i) in $store.state.groups" :key="group" :open="!i">
      <summary>Group: {{ group }}</summary>
      <GroupEntriesTable
        :group-id="group"
        :entries="$store.getters.entriesByGroup(group)"
      />
    </details>
    <div v-if="!$store.state.groups.length">No entries!!</div>
  </section>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { createGroupId, addCollectJob } from "../collect";
import { Config } from "../store";
import GroupEntriesTable from "./GroupEntriesTable.vue";

export default defineComponent({
  components: {
    GroupEntriesTable,
  },
  methods: {
    async collect() {
      const config = this.$store.getters.groupConfig as Config[];

      const GroupId = createGroupId();
      await Promise.all(
        config.map(({ Type, URL, Duration, Label }) =>
          addCollectJob(Type, { URL, Duration, Label, GroupId })
        )
      );
    },
  },
});
</script>

<style scoped lang="scss">
section {
  margin: 2em;
}

details {
  margin-bottom: 1em;
}

.control {
  margin-bottom: 2em;
}
</style>

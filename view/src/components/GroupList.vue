<template>
  <section>
    <div class="control">
      <button @click="collect">Collect</button>
      <router-link v-slot="{ navigate }" to="/group/config/" custom>
        <button @click="navigate">Configure</button>
      </router-link>
    </div>
    <details v-for="(group, i) in $store.state.groups" :key="group" :open="!i">
      <summary>Group: {{ group }}</summary>
      <GroupEntriesTable
        :group-id="group"
        :entries="$store.getters.entriesByGroup(group)"
      />
    </details>
  </section>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { createGroupId, addCollectJob } from "../collect";
import { Config, getGroupConfig } from "./GroupConfig.vue";
import GroupEntriesTable from "./GroupEntriesTable.vue";

export default defineComponent({
  components: {
    GroupEntriesTable,
  },
  methods: {
    async collect() {
      const config = JSON.parse(getGroupConfig()) as Config[];

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
  text-align: right;

  button {
    margin-left: 0.5em;
  }
}
</style>

<template>
  <div class="container">
    <nav>
      <router-link
        v-slot="{ navigate, isActive }"
        :to="`/group/${$route.params.gid}/index/`"
        custom
      >
        <div :class="{ active: isActive }" @click="navigate">index</div>
      </router-link>
      <router-link
        v-for="entry in $store.getters.availableEntriesByGroup(
          $route.params.gid
        )"
        v-slot="{ navigate, isActive }"
        :key="entry.Snapshot.ID"
        :to="`/group/${$route.params.gid}/${entry.Snapshot.Type}/${entry.Snapshot.ID}/`"
        custom
      >
        <div :class="{ active: isActive }" @click="navigate">
          {{ entry.Snapshot.Type }}: {{ entry.Snapshot.Label }}
        </div>
      </router-link>
    </nav>
    <router-view />
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { createGroupId, addCollectJob } from "../collect";
import { Config, getGroupConfig } from "./GroupConfig.vue";

export default defineComponent({
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
.container {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
}

nav {
  background-color: #555;
}
</style>

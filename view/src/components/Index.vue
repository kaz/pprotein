<template>
  <section>
    <div
      v-for="endpoint in $store.state.endpoints"
      :key="endpoint"
      class="container"
    >
      <h2>{{ endpoint }}</h2>
      <PproteinForm :ref="endpoint" :endpoint="endpoint" />
      <EntriesTable
        :entries="$store.state.entries[endpoint] || []"
        :length="4"
      />
    </div>
    <div class="control">
      <button @click="collect">Collect All</button>
    </div>
  </section>
</template>

<script lang="ts">
import EntriesTable from "./EntriesTable.vue";
import PproteinForm from "./PproteinForm.vue";
import { defineComponent } from "vue";

export default defineComponent({
  components: { EntriesTable, PproteinForm },
  methods: {
    collect() {
      return Promise.all(
        this.$store.state.endpoints.map((endpoint) =>
          (this.$refs[endpoint] as InstanceType<typeof PproteinForm>).collect()
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

.container {
  margin-bottom: 2em;
}

.control {
  margin: 3em 0;
  text-align: right;
}
</style>

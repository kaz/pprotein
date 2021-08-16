<template>
  <section>
    <div class="container right">
      <button @click="fetchAll">Fetch All</button>
    </div>
    <div v-for="key in Object.keys(eventSources)" :key="key" class="container">
      <h2>{{ key }}</h2>
      <PproteinForm :ref="key" :endpoint="key" />
      <EntriesTable
        :prefix="`/${key}`"
        :entries="$store.state.remote[key] || []"
        :length="3"
      />
    </div>
  </section>
</template>

<script lang="ts">
import EntriesTable from "./EntriesTable.vue";
import PproteinForm from "./PproteinForm.vue";
import { defineComponent } from "vue";

export default defineComponent({
  components: { EntriesTable, PproteinForm },
  data() {
    return {
      eventSources: {
        pprof: null,
        httplog: null,
        slowlog: null,
      } as { [key: string]: EventSource | null },
    };
  },
  beforeMount() {
    Object.keys(this.eventSources).forEach(this.subscribe);
  },
  beforeUnmount() {
    Object.keys(this.eventSources).forEach(this.unsubscribe);
  },
  methods: {
    fetchAll() {
      return Promise.all(
        Object.keys(this.eventSources).map((endpoint) =>
          (this.$refs[endpoint] as InstanceType<typeof PproteinForm>).fetch()
        )
      );
    },
    update(endpoint: string) {
      return this.$store.dispatch("syncStoreData", { endpoint });
    },
    subscribe(endpoint: string) {
      this.update(endpoint);

      const eventPath = `/api/${endpoint}/events`;
      const eventSource = new EventSource(eventPath);
      eventSource.onerror = () =>
        console.error(`EventSource error: ${eventPath}`);
      eventSource.onmessage = () => this.update(endpoint);
      this.eventSources[endpoint] = eventSource;
    },
    unsubscribe(endpoint: string) {
      this.eventSources[endpoint]?.close();
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

.right {
  text-align: right;
  button {
    margin-left: 1em;
  }
}
</style>

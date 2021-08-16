<template>
  <section>
    <PproteinForm :endpoint="$props.endpoint" @fetch="update" />
    <EntriesTable
      :endpoint="$props.endpoint"
      :entries="$store.state.remote[$props.endpoint]"
    />
  </section>
</template>

<script lang="ts">
import EntriesTable from "./EntriesTable.vue";
import PproteinForm from "./PproteinForm.vue";
import { defineComponent } from "vue";

export default defineComponent({
  components: {
    EntriesTable,
    PproteinForm,
  },
  props: {
    endpoint: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      eventSource: null as EventSource | null,
    };
  },
  beforeMount() {
    this.update();

    const eventPath = `/api/${this.$props.endpoint}/events`;
    this.eventSource = new EventSource(eventPath);
    this.eventSource.onerror = () =>
      console.error(`EventSource error: ${eventPath}`);
    this.eventSource.onmessage = () => this.update();
  },
  beforeUnmount() {
    this.eventSource?.close();
  },
  methods: {
    async update() {
      await this.$store.dispatch("syncStoreData", {
        endpoint: this.$props.endpoint,
      });
    },
  },
});
</script>

<style scoped lang="scss">
section {
  margin: 2em;
}
</style>

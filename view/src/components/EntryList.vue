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
import { defineComponent } from "vue";
import PproteinForm from "./PproteinForm.vue";
import EntriesTable from "./EntriesTable.vue";

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
      timer: -1,
    };
  },
  async beforeMount() {
    await this.update();
    this.$data.timer = setInterval(this.update, 2048);
  },
  async beforeUnmount() {
    clearInterval(this.$data.timer);
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

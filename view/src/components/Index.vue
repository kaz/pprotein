<template>
  <section>
    <div class="container">
      <h2>PProf</h2>
      <PproteinForm ref="pprof" endpoint="pprof/profiles" />
      <EntriesTable
        v-if="$store.state.remote['pprof/profiles']?.length > 0"
        prefix="/pprof"
        :entries="[$store.state.remote['pprof/profiles'][0]]"
      />
    </div>
    <div class="container">
      <h2>HTTP log</h2>
      <PproteinForm ref="httplog" endpoint="httplog/logs" />
      <EntriesTable
        v-if="$store.state.remote['httplog/logs']?.length > 0"
        prefix="/httplog"
        :entries="[$store.state.remote['httplog/logs'][0]]"
      />
    </div>
    <div class="container">
      <h2>MySQL slow log</h2>
      <PproteinForm ref="slowlog" endpoint="slowlog/logs" />
      <EntriesTable
        v-if="$store.state.remote['slowlog/logs']?.length > 0"
        prefix="/slowlog"
        :entries="[$store.state.remote['slowlog/logs'][0]]"
      />
    </div>
    <div class="container right">
      <button @click="updateState">Update State</button>
      <button @click="fetchAll">Fetch All</button>
    </div>
  </section>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import PproteinForm from "./PproteinForm.vue";
import EntriesTable from "./EntriesTable.vue";

export default defineComponent({
  components: { EntriesTable, PproteinForm },
  async beforeMount() {
    await this.updateState();
  },
  methods: {
    async fetchAll() {
      await Promise.all([
        (this.$refs.pprof as InstanceType<typeof PproteinForm>).fetch(),
        (this.$refs.httplog as InstanceType<typeof PproteinForm>).fetch(),
        (this.$refs.slowlog as InstanceType<typeof PproteinForm>).fetch(),
      ]);
      await this.updateState();
    },
    async updateState() {
      await Promise.all([
        this.$store.dispatch("syncStoreData", { endpoint: "pprof/profiles" }),
        this.$store.dispatch("syncStoreData", { endpoint: "httplog/logs" }),
        this.$store.dispatch("syncStoreData", { endpoint: "slowlog/logs" }),
      ]);
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

<template>
  <section>
    <table>
      <thead>
        <tr>
          <th>Open</th>
          <th>Datetime</th>
          <th>Source</th>
          <th>Duration</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
        <tr :key="info.Profile.ID" v-for="info in $store.state.profiles">
          <td>
            <router-link :to="`/pprof/${info.Profile.ID}`">Open</router-link>
          </td>
          <td>{{ info.Profile.Datetime }}</td>
          <td>{{ info.Profile.URL }}</td>
          <td>{{ info.Profile.Duration }}</td>
          <td>{{ info.Status }}</td>
        </tr>
      </tbody>
    </table>
  </section>
</template>

<style scoped lang="scss">
section {
  margin: 2em;
}

table {
  border-collapse: collapse;
}
th,
td {
  padding: 0.5em 2em;
  border: 1px solid #999;
}
</style>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  data() {
    return {
      timer: -1,
    };
  },
  async beforeMount() {
    const updateFn = () => this.$store.dispatch("updateProfiles");
    await updateFn();
    this.$data.timer = setInterval(updateFn, 2048);
  },
  async beforeUnmount() {
    clearInterval(this.$data.timer);
  },
});
</script>

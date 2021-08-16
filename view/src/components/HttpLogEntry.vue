<template>
  <pre>{{ $data.summary }}</pre>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  data() {
    return {
      summary: "Loading ...",
    };
  },
  async beforeCreate() {
    const resp = await fetch(`/api/httplog/logs/${this.$route.params.id}`);
    this.$data.summary = await resp.text();
  },
});
</script>

<style scoped lang="scss">
pre {
  margin: 0;
  padding: 2em;
  overflow: auto;
  width: 100vw;
  flex: 1 0 auto;
}
</style>

<template>
  <section>
    <pre>{{ $data.memo.Text || $data.summary }}</pre>
  </section>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  data() {
    return {
      summary: "Loading ...",
      memo: { Text: null },
    };
  },
  async beforeCreate() {
    try {
      const resp = await fetch(`/api/memo/${this.$route.params.id}`);
      this.memo = await resp.json();
    } catch (e) {
      this.summary = `Error: ${e instanceof Error ? e.message : e}`;
    }
  },
});
</script>
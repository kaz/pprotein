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
  async created() {
    await this.updateMemo(this.$route.params.id);
  },
  async beforeRouteUpdate(route) {
    await this.updateMemo(route.params.id);
  },
  methods: {
    async updateMemo(id: string | string[]) {
      try {
        const resp = await fetch(`/api/memo/${id}`);
        this.memo = await resp.json();
      } catch (e) {
        this.summary = `Error: ${e instanceof Error ? e.message : e}`;
      }
    },
  },
});
</script>

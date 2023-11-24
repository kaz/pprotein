<template>
  <TsvTable :tsv="tsv" :link="`/api/httplog/data/${$route.params.id}`"/>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import TsvTable from "./TsvTable.vue";

export default defineComponent({
  components: {
    TsvTable,
  },
  data() {
    return {
      tsv: "",
    };
  },
  async created() {
    await this.updateTsv(this.$route.params.id);
  },
  async beforeRouteUpdate(route) {
    await this.updateTsv(route.params.id);
  },
  methods: {
    async updateTsv(id: string | string[]) {
      const resp = await fetch(`/api/httplog/${id}`);
      this.tsv = await resp.text();
    },
  },
});
</script>

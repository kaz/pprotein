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
  async beforeCreate() {
    const resp = await fetch(`/api/httplog/${this.$route.params.id}`);
    this.tsv = await resp.text();
  },
});
</script>

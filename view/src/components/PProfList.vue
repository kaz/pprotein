<template>
  <EntryList endpoint="pprof" />
</template>

<script lang="ts">
import EntryList from "./EntryList.vue";
import { defineComponent } from "vue";

export default defineComponent({
  components: {
    EntryList,
  },
});
</script>

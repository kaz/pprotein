<template>
  <EntryList endpoint="pprof/profiles" />
</template>

<script lang="ts">
import { defineComponent } from "vue";
import EntryList from "./EntryList.vue";

export default defineComponent({
  components: {
    EntryList,
  },
});
</script>

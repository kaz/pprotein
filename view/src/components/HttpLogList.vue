<template>
  <EntryList endpoint="httplog" />
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

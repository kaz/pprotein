<template>
  <iframe :src="`/api/pprof/${$route.params.id}/flamegraph`" />
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({});
</script>

<style scoped lang="scss">
iframe {
  flex: 1 0 auto;
}
</style>

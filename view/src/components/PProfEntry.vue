<template>
  <iframe :src="`/api/pprof/profiles/${$route.params.id}/`" />
</template>

<style scoped lang="scss">
iframe {
  width: 100vw;
  flex: 1 0 auto;
}
</style>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({});
</script>

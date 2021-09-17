<template>
  <span
    v-if="
      $props.sha in $store.state.commits &&
      $store.state.commits[$props.sha].html_url
    "
  >
    <a target="_blank" :href="$store.state.commits[$props.sha].html_url"
      >[{{ $props.sha.substr(0, 7) }}]</a
    >
    {{ $store.state.commits[$props.sha].message }}
    ({{ $store.state.commits[$props.sha].author.name }})
  </span>
  <span v-else>[unknown]</span>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  props: {
    sha: {
      type: String,
      required: true,
    },
  },
});
</script>

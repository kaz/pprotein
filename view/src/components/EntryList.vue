<template>
  <section>
    <PproteinForm :endpoint="$props.endpoint" />
    <EntriesTable
      :endpoint="$props.endpoint"
      :entries="$store.state.entries[$props.endpoint]"
    />
  </section>
</template>

<script lang="ts">
import EntriesTable from "./EntriesTable.vue";
import PproteinForm from "./PproteinForm.vue";
import { defineComponent } from "vue";

export default defineComponent({
  components: {
    EntriesTable,
    PproteinForm,
  },
  props: {
    endpoint: {
      type: String,
      required: true,
    },
  },
});
</script>

<style scoped lang="scss">
section {
  margin: 2em;
}
</style>

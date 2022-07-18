<template>
  <table>
    <thead>
      <tr>
        <th></th>
        <th>Datetime</th>
        <th>Source URL</th>
        <th>Duration</th>
        <th>Commit</th>
        <th>Status</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="entry in visibleEntries" :key="entry.Snapshot.ID">
        <td>
          <router-link
            v-if="entry.Status == `ok`"
            :to="
              entry.Snapshot.GroupId
                ? `/group/${entry.Snapshot.GroupId}/${entry.Snapshot.Type}/${entry.Snapshot.ID}/`
                : `/${entry.Snapshot.Type}/${entry.Snapshot.ID}/`
            "
          >
            Open
          </router-link>
        </td>
        <td>{{ entry.Snapshot.Datetime.toLocaleString() }}</td>
        <td>{{ entry.Snapshot.URL }}</td>
        <td>{{ entry.Snapshot.Duration }}</td>
        <td><Commit :repository="entry.Snapshot.Repository" /></td>
        <td><Status :status="entry.Status" :message="entry.Message" /></td>
      </tr>
    </tbody>
  </table>
</template>

<script lang="ts">
import { defineComponent, PropType } from "vue";
import { Entry } from "../store";
import Commit from "./Commit.vue";
import Status from "./Status.vue";

export default defineComponent({
  components: {
    Commit,
    Status,
  },
  props: {
    entries: {
      type: Array as PropType<Entry[]>,
      required: true,
    },
    length: {
      type: Number,
      default: undefined,
    },
  },
  computed: {
    visibleEntries() {
      return this.$props.length
        ? this.$props.entries.slice(0, this.$props.length)
        : this.$props.entries;
    },
  },
});
</script>

<style scoped lang="scss">
table {
  border-collapse: collapse;
}

th,
td {
  padding: 0.5em 2em;
  border: 1px solid #999;
}
</style>

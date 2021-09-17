<template>
  <table>
    <thead>
      <tr>
        <th>Open</th>
        <th>Type</th>
        <th>Label</th>
        <th>Commit</th>
        <th>Status</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="entry in visibleEntries" :key="entry.Snapshot.ID">
        <td>
          <router-link
            v-if="entry.Status == `ok`"
            :to="`/group/${$props.groupId}/${entry.Snapshot.Type}/${entry.Snapshot.ID}/`"
            >Open</router-link
          >
        </td>
        <td>{{ entry.Snapshot.Type }}</td>
        <td>{{ entry.Snapshot.Label }}</td>
        <td><Commit :sha="entry.Snapshot.GitRevision" /></td>
        <td>
          <div :class="`cell ${entry.Status}`"></div>
          {{ entry.Message || entry.Status }}
        </td>
      </tr>
    </tbody>
  </table>
</template>

<script lang="ts">
import { Entry } from "../store";
import { PropType, defineComponent } from "vue";
import Commit from "./Commit.vue";

export default defineComponent({
  components: {
    Commit,
  },
  props: {
    groupId: {
      type: String,
      required: true,
    },
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

.cell {
  display: inline-block;
  width: 1em;
  height: 1em;
  border-radius: 0.2em;
  top: 0.1em;
  position: relative;

  &.ok {
    background-color: blue;
  }
  &.fail {
    background-color: red;
  }
  &.pending {
    background-color: orange;
    animation: flash 1s ease-in-out 0s infinite alternate;
  }
}

@keyframes flash {
  0% {
    opacity: 1;
  }
  100% {
    opacity: 0.1;
  }
}
</style>

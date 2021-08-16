<template>
  <table>
    <thead>
    <tr>
      <th>Open</th>
      <th>Datetime</th>
      <th>Source URL</th>
      <th>Duration</th>
      <th>Status</th>
    </tr>
    </thead>
    <tbody>
    <tr :key="info.Entry.ID" v-for="info in $props.entries">
      <td>
        <router-link v-if="info.Status == `ok`" :to="`${$props.prefix ? $props.prefix : '.'}/${info.Entry.ID}/`">Open</router-link>
      </td>
      <td>{{ info.Entry.Datetime.toLocaleString() }}</td>
      <td>{{ info.Entry.URL }}</td>
      <td>{{ info.Entry.Duration }}</td>
      <td>
        <div :class="`cell ${info.Status}`"></div>
        {{ info.Message || info.Status }}
      </td>
    </tr>
    </tbody>
  </table>
</template>

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

<script lang="ts">
export default {
  name: 'EntriesTable',
  props: {
    prefix: String,
    entries: Array,
  },
}
</script>

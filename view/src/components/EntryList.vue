<template>
  <section>
    <div class="form">
      <label>
        Source URL<br />
        <input type="text" size="50" v-model="$data.url" />
      </label>
      <label>
        Duration<br />
        <input type="number" size="10" v-model.number="$data.duration" />
      </label>
      <label>
        &nbsp;<br />
        <button @click="fetch">Fetch</button>
      </label>
    </div>
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
        <tr :key="info.Entry.ID" v-for="info in $store.state.remote[$props.endpoint]">
          <td>
            <router-link v-if="info.Status == `ok`" :to="`./${info.Entry.ID}/`">Open</router-link>
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
  </section>
</template>

<style scoped lang="scss">
section {
  margin: 2em;
}

.form {
  display: flex;
  margin-bottom: 2em;
  label {
    margin-right: 1em;
    &:focus-within {
      color: orangered;
    }
    input {
      border: 1px solid lightgray;
      padding: 0.4em 1em;
      &:focus {
        border-color: orangered;
        outline: 0;
      }
    }
    button {
      padding: 0.4em 1em;
      background-color: white;
      border: 1px solid lightgray;
      cursor: pointer;
      &:hover {
        border-color: orangered;
      }
      &:active {
        color: white;
        background-color: orangered;
      }
    }
  }
}

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
import { defineComponent } from "vue";

export default defineComponent({
  props: {
    endpoint: String,
  },
  data() {
    return {
      url: localStorage.getItem(`url[${this.$props.endpoint}]`) || "http://",
      duration: localStorage.getItem(`duration[${this.$props.endpoint}]`) || 60,
      timer: -1,
    };
  },
  watch: {
    url(val) {
      localStorage.setItem(`url[${this.$props.endpoint}]`, val);
    },
    duration(val) {
      localStorage.setItem(`duration[${this.$props.endpoint}]`, val);
    },
  },
  methods: {
    async update() {
      await this.$store.dispatch("syncStoreData", { endpoint: this.$props.endpoint });
    },
    async fetch() {
      await this.$store.dispatch("postStoreData", {
        endpoint: this.$props.endpoint,
        URL: this.$data.url,
        Duration: parseInt(this.$data.duration),
      });
      await this.update();
    },
  },
  async beforeMount() {
    await this.update();
    this.$data.timer = setInterval(this.update, 2048);
  },
  async beforeUnmount() {
    clearInterval(this.$data.timer);
  },
});
</script>

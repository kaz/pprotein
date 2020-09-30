<template>
  <section>
    <div class="form">
      <label>
        Source URL<br />
        <input v-model="url" type="text" size="45" />
      </label>
      <label>
        Duration<br />
        <input v-model.number="duration" type="number" />
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
        <tr :key="info.Profile.ID" v-for="info in $store.state.profiles">
          <td>
            <router-link v-if="info.Status == `ok`" :to="`/pprof/${info.Profile.ID}`">Open</router-link>
          </td>
          <td>{{ info.Profile.Datetime.toLocaleString() }}</td>
          <td>{{ info.Profile.URL }}</td>
          <td>{{ info.Profile.Duration }}</td>
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
  data() {
    return {
      timer: -1,
      url: "",
      duration: 60,
    };
  },
  methods: {
    async update() {
      await this.$store.dispatch("updateProfiles");
    },
    async fetch() {
      localStorage.setItem("saved_url", this.$data.url);
      localStorage.setItem("saved_duration", this.$data.duration);

      const resp = await fetch("/api/pprof/profiles", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          URL: this.$data.url,
          Duration: this.$data.duration,
        }),
      });

      if (!resp.ok) {
        alert(await resp.text());
      }
      await this.update();
    },
  },
  async beforeMount() {
    this.$data.url = localStorage.getItem("saved_url") || "";
    this.$data.duration = parseInt(localStorage.getItem("saved_duration")) || 60;

    await this.update();
    this.$data.timer = setInterval(this.update, 2048);
  },
  async beforeUnmount() {
    clearInterval(this.$data.timer);
  },
});
</script>

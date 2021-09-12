<template>
  <section>
    <textarea v-model="$data.content" :disabled="processing"></textarea>
    <div class="control">
      <button :disabled="processing" @click="update">
        {{ $data.label || "Update" }}
      </button>
    </div>
  </section>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { SnapshotTarget } from "../store";

export type Config = { Type: string } & Omit<SnapshotTarget, "GroupId">;
export const localStorageKey = "groupConfig";

const defaultConfig: Config[] = [
  {
    Type: "pprof",
    URL: "http://localhost:9000/debug/pprof/profile",
    Duration: 10,
    Label: "101-localhost",
  },
  {
    Type: "httplog",
    URL: "http://localhost:9000/debug/log/httplog",
    Duration: 10,
    Label: "102-localhost",
  },
  {
    Type: "slowlog",
    URL: "http://localhost:9000/debug/log/slowlog",
    Duration: 10,
    Label: "103-localhost",
  },
];

export const getGroupConfig = (): string =>
  localStorage.getItem(localStorageKey) ||
  JSON.stringify(defaultConfig, null, 4);

export default defineComponent({
  data() {
    return {
      processing: false,
      label: "",
      content: getGroupConfig(),
    };
  },
  methods: {
    async update() {
      this.$data.processing = true;
      this.$data.label = "Updating ...";

      try {
        const contentObj = JSON.parse(this.$data.content);
        localStorage.setItem(
          localStorageKey,
          JSON.stringify(contentObj, null, 4)
        );
      } catch (e) {
        alert(`ERROR: ${e}`);

        this.$data.label = "Failed";
        this.reset();
        return;
      }

      this.$data.label = "Updated!";
      this.reset();
    },
    reset() {
      setTimeout(() => {
        this.$data.processing = false;
        this.$data.label = "";
      }, 2000);
    },
  },
});
</script>

<style scoped lang="scss">
section {
  margin: 2em;
}

textarea {
  width: 100%;
  height: 70vh;
}

.control {
  margin: 3em 0;
  text-align: right;
}
</style>

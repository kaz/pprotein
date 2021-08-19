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

export default defineComponent({
  data() {
    return {
      processing: true,
      label: "",
      content: "Loading ...",
    };
  },
  async beforeCreate() {
    const resp = await fetch(`/api/httplog/config`);
    this.$data.content = await resp.text();
    this.$data.processing = false;
  },
  methods: {
    async update() {
      this.$data.processing = true;
      this.$data.label = "Updating ...";

      const resp = await fetch(`/api/httplog/config`, {
        method: "POST",
        body: this.$data.content,
      });
      if (!resp.ok) {
        alert(`ERROR: ${await resp.text()}`);
        return;
      }

      this.$data.label = "Updated!";
      setTimeout(() => {
        this.$data.processing = false;
        this.$data.label = "";
      }, 2048);
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

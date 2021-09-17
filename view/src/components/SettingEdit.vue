<template>
  <div>
    <textarea v-model="$data.content" :disabled="!!$data.status"></textarea>
    <div class="control">
      <button :disabled="!!$data.status" @click="update">
        {{ $data.status || "Update" }}
      </button>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  props: {
    name: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      status: "",
      content: "Loading ...",
    };
  },
  beforeMount() {
    this.sync();
    this.$store.watch(
      () => this.$store.state.settings[this.$props.name],
      this.sync
    );
  },
  methods: {
    sync() {
      const ent = this.$store.state.settings[this.$props.name];
      if (ent) {
        this.$data.content = ent.value;
      }
    },
    async update() {
      this.$data.status = "Updating ...";

      await this.$store.dispatch("updateSetting", {
        key: this.$props.name,
        value: this.$data.content,
      });

      this.$data.status = "Updated!";
      setTimeout(() => (this.$data.status = ""), 2000);
    },
  },
});
</script>

<style scoped lang="scss">
textarea {
  width: 100%;
  height: 25vh;
}

.control {
  margin: 1em 0;
  text-align: right;
}
</style>

<template>
  <section>
    <template v-for="key in $store.state.settingKeys" :key="key">
      <h1>{{ key }}</h1>
      <SettingEdit :name="key" />
    </template>
  </section>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import SettingEdit from "./SettingEdit.vue";

export default defineComponent({
  components: { SettingEdit },
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
</style>

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
});
</script>

<style scoped lang="scss">
section {
  margin: 2em;
}
</style>

<template>
  <section>
    <div class="control">
      <button>Collect</button>
      <router-link v-slot="{ navigate }" to="/all/config/" custom>
        <button @click="navigate">Configure</button>
      </router-link>
    </div>
  </section>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({});
</script>

<style scoped lang="scss">
section {
  margin: 2em;
}

.control {
  text-align: right;

  button {
    margin-left: 0.5em;
  }
}
</style>

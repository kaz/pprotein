<template>
  <div>
    <EntryList endpoint="httplog" />
    <div class="control">
      <router-link v-slot="{ navigate }" to="/httplog/config/" custom>
        <button @click="navigate">Edit kataribe.toml</button>
      </router-link>
    </div>
  </div>
</template>

<script lang="ts">
import EntryList from "./EntryList.vue";
import { defineComponent } from "vue";

export default defineComponent({
  components: {
    EntryList,
  },
});
</script>
<style scoped lang="scss">
.control {
  margin: 3em 4em;
  text-align: right;
}
</style>

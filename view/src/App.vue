<template>
  <main>
    <header>PProtein ⚙ Manage Panel</header>
    <nav>
      <router-link to="/pprof/" custom v-slot="{ navigate, isActive }">
        <div :class="{ active: isActive }" @click="navigate">PProf</div>
      </router-link>
      <router-link to="/httplog/" custom v-slot="{ navigate, isActive }">
        <div :class="{ active: isActive }" @click="navigate">HTTP Log</div>
      </router-link>
      <router-link to="/slowlog/" custom v-slot="{ navigate, isActive }">
        <div :class="{ active: isActive }" @click="navigate">Slow Log</div>
      </router-link>
    </nav>
    <router-view />
  </main>
</template>

<style lang="scss">
* {
  font-size: 14px;
  font-family: "Courier Prime", monospace;
  box-sizing: border-box;
}
body,
html {
  padding: 0;
  margin: 0;
}
main {
  display: flex;
  flex-direction: column;
}
header {
  padding: 1em 2em;
  background-color: #111;
  color: #fff;
}
nav {
  display: flex;
  background-color: #333;
  color: #fff;
  div {
    cursor: pointer;
    padding: 0.7em 2em 0.4em 2em;
    border-bottom: 0.3em solid transparent;
    &.active {
      border-bottom: 0.3em solid orange;
    }
  }
}
section {
  padding: 1em 2em;
}
</style>

<script lang="ts">
import { defineComponent } from "vue";

import "typeface-courier-prime";

type Dict = { [key: string]: string };

export default defineComponent({
  methods: {
    getTitle(params: Dict, meta: Dict) {
      return Object.entries(params).reduce((title, [key, val]) => title.replace(`{{${key}}}`, val), meta.title || "");
    },
  },
  watch: {
    $route({ params, meta }) {
      document.title = `${this.getTitle(params, meta)} | PProtein`;
    },
  },
});
</script>

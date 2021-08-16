<template>
  <main>
    <header>PProtein ⚙ Manage Panel</header>
    <nav>
      <router-link v-slot="{ navigate, isActive }" to="/" custom>
        <div :class="{ active: isActive }" @click="navigate">Top</div>
      </router-link>
      <router-link v-slot="{ navigate, isActive }" to="/pprof/" custom>
        <div :class="{ active: isActive }" @click="navigate">PProf</div>
      </router-link>
      <router-link v-slot="{ navigate, isActive }" to="/httplog/" custom>
        <div :class="{ active: isActive }" @click="navigate">HTTP Log</div>
      </router-link>
      <router-link v-slot="{ navigate, isActive }" to="/slowlog/" custom>
        <div :class="{ active: isActive }" @click="navigate">Slow Log</div>
      </router-link>
    </nav>
    <router-view />
  </main>
</template>

<script lang="ts">
import { defineComponent } from "vue";

import "typeface-courier-prime";

type Dict = { [key: string]: string };

export default defineComponent({
  watch: {
    $route({ params, meta }) {
      document.title = `${this.getTitle(params, meta)} | PProtein`;
    },
  },
  methods: {
    getTitle(params: Dict, meta: Dict) {
      return Object.entries(params).reduce(
        (title, [key, val]) => title.replace(`{{${key}}}`, val),
        meta.title || ""
      );
    },
  },
});
</script>

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
</style>

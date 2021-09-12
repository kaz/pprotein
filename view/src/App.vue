<template>
  <main>
    <header>
      <router-link to="/">{{ $data.title }}</router-link>
    </header>
    <nav>
      <router-link v-slot="{ navigate, isActive }" to="/group/" custom>
        <div :class="{ active: isActive }" @click="navigate">group</div>
      </router-link>
      <router-link v-slot="{ navigate, isActive }" to="/pprof/" custom>
        <div :class="{ active: isActive }" @click="navigate">pprof</div>
      </router-link>
      <router-link v-slot="{ navigate, isActive }" to="/httplog/" custom>
        <div :class="{ active: isActive }" @click="navigate">httplog</div>
      </router-link>
      <router-link v-slot="{ navigate, isActive }" to="/slowlog/" custom>
        <div :class="{ active: isActive }" @click="navigate">slowlog</div>
      </router-link>
    </nav>
    <router-view />
  </main>
</template>

<script lang="ts">
import "typeface-courier-prime";
import { defineComponent } from "vue";

type Dict = { [key: string]: string };

export default defineComponent({
  data() {
    return {
      title: "pprotein ⚙",
    };
  },
  watch: {
    $route({ params, meta }) {
      document.title = `${this.getTitle(params, meta)} | ${this.$data.title}`;
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
  height: 100vh;
  width: 100vw;
}

header {
  flex-shrink: 0;
  padding: 1em 2em;
  background-color: #111;

  a {
    color: #fff;
    text-decoration: none;
  }
}

nav {
  flex-shrink: 0;
  display: flex;
  overflow: auto;
  background-color: #333;
  color: #fff;

  div {
    cursor: pointer;
    white-space: nowrap;
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

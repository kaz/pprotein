<template>
  <div class="form">
    <label>
      Source URL<br/>
      <input type="text" size="50" v-model="$data.url"/>
    </label>
    <label>
      Duration<br/>
      <input type="number" size="10" v-model.number="$data.duration"/>
    </label>
    <label>
      &nbsp;<br/>
      <button @click="fetch">Fetch</button>
    </label>
  </div>
</template>

<style scoped lang="scss">
.form {
  display: flex;
  margin-bottom: 2em;

  label {
    margin-right: 1em;

    &:focus-within {
      color: orangered;
    }

    input {
      border: 1px solid lightgray;
      padding: 0.4em 1em;

      &:focus {
        border-color: orangered;
        outline: 0;
      }
    }

  }
}
</style>

<script lang="ts">
export default {
  name: 'PproteinForm',
  props: {
    endpoint: String
  },
  data() {
    return {
      url: localStorage.getItem(`url[${this.$props.endpoint}]`) || "http://",
      duration: localStorage.getItem(`duration[${this.$props.endpoint}]`) || 60,
    };
  },
  methods: {
    async fetch() {
      await this.$store.dispatch("postStoreData", {
        endpoint: this.$props.endpoint,
        URL: this.$data.url,
        Duration: parseInt(this.$data.duration),
      });
      this.$emit("fetch");
    },
  },
  watch: {
    url(val) {
      localStorage.setItem(`url[${this.$props.endpoint}]`, val);
    },
    duration(val) {
      localStorage.setItem(`duration[${this.$props.endpoint}]`, val);
    },
  },
}
</script>
<template>
  <div class="form">
    <label>
      Source URL<br />
      <input v-model="$data.url" type="text" size="50" />
    </label>
    <label>
      Duration<br />
      <input v-model.number="$data.duration" type="number" size="10" />
    </label>
    <label>
      &nbsp;<br />
      <button @click="collect">Collect</button>
    </label>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { addCollectJob } from "../collect";

export default defineComponent({
  props: {
    endpoint: {
      type: String,
      required: true,
    },
  },
  data(): { url: string; duration: string } {
    return {
      url: localStorage.getItem(`url[${this.$props.endpoint}]`) || "http://",
      duration:
        localStorage.getItem(`duration[${this.$props.endpoint}]`) || "60",
    };
  },
  watch: {
    url(val) {
      localStorage.setItem(`url[${this.$props.endpoint}]`, val);
    },
    duration(val) {
      localStorage.setItem(`duration[${this.$props.endpoint}]`, val);
    },
  },
  methods: {
    async collect() {
      await addCollectJob(this.$props.endpoint, {
        GroupId: "",
        Label: "",
        URL: this.$data.url,
        Duration: parseInt(this.$data.duration),
      });
    },
  },
});
</script>

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

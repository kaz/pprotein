<template>
  <div class="add-memo-container">
    <div>
      <label>
        Label:
        <input v-model="label" type="text" />
      </label>
      <div style="margin-top: 0.4rem">
        <label>
          Text :
          <textarea v-model="text" cols="60" rows="3"></textarea>
        </label>
      </div>
    </div>
    <button @click="addComment">Add Memo</button>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  props: {
    groupId: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      label: "",
      text: "",
    };
  },
  methods: {
    async addComment() {
      const resp = await fetch(`/api/memo`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          GroupId: this.groupId,
          Text: this.text,
          Label: this.label,
        }),
      });

      if (!resp.ok) {
        alert(await resp.text());
      }
    },
  },
});
</script>

<style scoped lang="scss">
.add-memo-container {
  margin-top: 1em;
  display: flex;
  align-items: end;
  gap: 1em;
}
</style>

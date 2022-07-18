<template>
  <div class="wrap">
    <div :class="['indicator', $props.status]" />
    <a
      v-if="$props.status == 'fail' && !openDetail"
      @click="showDetail"
      href="javascript:"
    >
      Failed …
    </a>
    <span v-else>
      {{ $props.message || $props.status }}
    </span>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from "vue";
import { StatusText } from "../store";

export default defineComponent({
  props: {
    status: {
      type: String as PropType<StatusText>,
      required: true,
    },
    message: {
      type: String,
    },
  },
  data: () => ({
    openDetail: false,
  }),
  methods: {
    showDetail() {
      this.openDetail = true;
    },
  },
});
</script>

<style scoped lang="scss">
.wrap {
  display: flex;
}

@keyframes flash {
  0% {
    opacity: 1;
  }
  100% {
    opacity: 0.1;
  }
}

.indicator {
  flex: 0 0 auto;
  margin-right: 0.4em;

  width: 1em;
  height: 1em;
  border-radius: 0.2em;

  &.ok {
    background-color: blue;
  }
  &.fail {
    background-color: red;
  }
  &.pending {
    background-color: orange;
    animation: flash 1s ease-in-out 0s infinite alternate;
  }
}
</style>

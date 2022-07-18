<template>
  <span v-if="$props.repository">
    <a v-if="!openDetail" @click="showDetail" href="javascript:">
      {{ $props.repository.Message.split("\n")[0] }} …
    </a>
    <template v-else>
      <span>
        {{ $props.repository.Message }}
      </span>
      <br />
      <small>
        {{ $props.repository.Author }}
        <br />
        <a target="_blank" :href="commitUrl">
          {{ $props.repository.Hash.substring(0, 7) }}
        </a>
        {{ " " }}
        <a target="_blank" :href="treeUrl">
          {{ $props.repository.Ref }}
        </a>
      </small>
    </template>
  </span>
  <span v-else>[unknown]</span>
</template>

<script lang="ts">
import { defineComponent, PropType } from "vue";
import { RepositoryInfo } from "../store";

export default defineComponent({
  props: {
    repository: {
      type: Object as PropType<RepositoryInfo>,
    },
  },
  data: () => ({
    openDetail: false,
  }),
  computed: {
    repoUrl() {
      const remoteUrl = this.$props.repository?.Remote;
      if (!remoteUrl) {
        return;
      }

      const sshMatch = remoteUrl.match(/^git@github\.com:(.+?)\/(.+?)\.git$/);
      if (sshMatch) {
        return `https://github.com/${sshMatch[1]}/${sshMatch[2]}`;
      }

      const httpsMatch = remoteUrl.match(
        /^https:\/\/github\.com\/(.+?)\/(.+?)\.git$/
      );
      if (httpsMatch) {
        return `https://github.com/${httpsMatch[1]}/${httpsMatch[2]}`;
      }

      return;
    },
    commitUrl() {
      const repoUrl = this.repoUrl;
      return repoUrl
        ? `${repoUrl}/commit/${this.$props.repository?.Hash}`
        : undefined;
    },
    treeUrl() {
      const [, branch] =
        this.$props.repository?.Ref.match(/^refs\/heads\/(.+)$/) || [];

      const repoUrl = this.repoUrl;
      return repoUrl && branch ? `${repoUrl}/tree/${branch}` : undefined;
    },
  },
  methods: {
    showDetail() {
      this.openDetail = true;
    },
  },
});
</script>

<style scoped lang="scss">
small {
  display: inline-block;
  margin-top: 0.8em;
  font-size: 0.8em;
  line-height: 0.8em;
}
</style>

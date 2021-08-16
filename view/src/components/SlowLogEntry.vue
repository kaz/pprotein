<template>
  <section>
    <pre v-if="!slowData.global">{{ $data.summary }}</pre>
    <div v-else class="container">
      <div>
        <h3>全体統計</h3>
        <div>クエリ数: {{ slowData?.global?.query_count }}</div>
        <div>ユニーククエリ数: {{ slowData?.global?.unique_query_count }}</div>
        <table border="1">
          <thead>
            <tr>
              <th>項目</th>
              <th>avg</th>
              <th>max</th>
              <th>median</th>
              <th>min</th>
              <th>pct_95</th>
              <th>stddev</th>
              <th>sum</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(row, key) in slowData?.global?.metrics" :key="key">
              <td>{{ key }}</td>
              <td v-for="(val, idx) in Object.values(row)" :key="idx">
                {{ val }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div>
        <h3>クエリ統計</h3>
        <div>
          <label
            ><input
              id="total_time"
              v-bind="inputAttrs('total_time')"
              type="radio"
              value="total_time"
            />合計時間</label
          >
          <label
            ><input
              id="count"
              v-bind="inputAttrs('count')"
              type="radio"
              value="count"
            />クエリ回数</label
          >
          <label
            ><input
              id="avg_time"
              v-bind="inputAttrs('avg_time')"
              type="radio"
              value="avg_time"
            />平均時間</label
          >
          <label
            ><input
              id="rows_s/e"
              v-bind="inputAttrs('rows_s/e')"
              type="radio"
              value="rows_s/e"
            />行効率</label
          >
        </div>
        <div v-for="queryClass in classes" :key="queryClass.checksum">
          <details>
            <summary>{{ queryClass.fingerprint }}</summary>
            <div>
              <div>クエリ数: {{ queryClass.query_count }}</div>
              <button @click="copy(queryClass.example.query)">
                サンプルクエリのコピー
              </button>
              <table border="1">
                <thead>
                  <tr>
                    <th>項目</th>
                    <th v-for="key in queryMetricsColKeys" :key="key">
                      {{ key }}
                    </th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="key in queryMetricsRowKeys" :key="key">
                    <td>{{ key }}</td>
                    <td v-for="colKey in queryMetricsColKeys" :key="colKey">
                      {{ queryClass.metrics[key][colKey] }}
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </details>
        </div>
      </div>
    </div>
  </section>
</template>

<script lang="ts">
import { Class, Metrics, MetricsRow, QueryDigest } from "../query-digest";
import { defineComponent } from "vue";

export default defineComponent({
  data() {
    return {
      summary: "Loading ...",
      slowData: {} as QueryDigest,
      sort: "total_time",
      queryMetricsRowKeys: [
        "Lock_time",
        "Query_length",
        "Query_time",
        "Rows_examined",
        "Rows_sent",
      ] as Array<keyof Metrics>,
      queryMetricsColKeys: [
        "sum",
        "avg",
        "max",
        "median",
        "min",
        "pct",
        "pct_95",
        "stddev",
      ] as Array<keyof MetricsRow>,
    };
  },
  computed: {
    classes() {
      if (!this.slowData) return [];
      const arr: Class[] = Array.from(this.slowData.classes);
      switch (this.sort) {
        case "count":
          return arr.sort((a, b) => b.query_count - a.query_count);
        case "total_time":
          return arr.sort(
            (a, b) =>
              Number(b.metrics.Query_time.sum) -
              Number(a.metrics.Query_time.sum)
          );
        case "avg_time":
          return arr.sort(
            (a, b) =>
              Number(b.metrics.Query_time.avg) -
              Number(a.metrics.Query_time.avg)
          );
        case "rows_s/e":
          return arr.sort(
            (a, b) =>
              Number(a.metrics.Rows_sent.sum) /
                Number(a.metrics.Rows_examined.sum) -
              Number(b.metrics.Rows_sent.sum) /
                Number(b.metrics.Rows_examined.sum)
          );
        default:
          return arr;
      }
    },
  },
  async beforeCreate() {
    const resp = await fetch(`/api/slowlog/${this.$route.params.id}`);
    this.slowData = await resp.json();
  },
  methods: {
    test() {
      console.log(this.$data.slowData.classes);
    },
    copy(text: string) {
      window.navigator.clipboard.writeText(text);
    },
    // このワークアラウンドを利用 https://github.com/johnsoncodehk/volar/issues/369#issuecomment-898250309
    inputAttrs(value: string) {
      return {
        value,
        checked: this.sort === value,
        onInput: (e: Event) => {
          const el = e.target as HTMLInputElement;
          this.sort = el.value;
        },
      };
    },
  },
});
</script>

<style scoped lang="scss">
.container {
  padding: 16px;
}

table {
  margin: 8px 0;
  border-collapse: collapse;
}

td,
th {
  padding: 4px;
}

details {
  border-bottom: solid 1px #111111;
  & > div {
    padding: 12px;
  }
}

summary {
  font-weight: bold;
  background-color: #dcdcdc;
  padding: 8px;
  overflow: hidden;
}
</style>

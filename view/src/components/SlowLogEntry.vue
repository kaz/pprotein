<template>
  <div>
    <pre v-if="!slowData">{{ $data.summary }}</pre>
    <div class="container" v-else>
      <div>
        <h3>全体統計</h3>
        <div>クエリ数: {{slowData.global.query_count}}</div>
        <div>ユニーククエリ数: {{slowData.global.unique_query_count}}</div>
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
            <tr v-for="(row, key) in slowData.global.metrics" :key="key">
              <td>{{key}}</td>
              <td v-for="val in Object.values(row)">{{val}}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <div>
        <h3>クエリ統計</h3>
        <div>
          <label><input type="radio" id="total_time" value="total_time" v-model="sort">合計時間</label>
          <label><input type="radio" id="count" value="count" v-model="sort">クエリ回数</label>
          <label><input type="radio" id="avg_time" value="avg_time" v-model="sort">平均時間</label>
          <label><input type="radio" id="rows_s/e" value="rows_s/e" v-model="sort">行効率</label>
        </div>
        <div v-for="queryClass in classes" :key="queryClass.checksum">
          <details>
            <summary>{{queryClass.fingerprint}}</summary>
            <div>
              <div>クエリ数: {{queryClass.query_count}}</div>
              <button @click="copy(queryClass.example.query)">サンプルクエリのコピー</button>
              <table border="1">
                <thead>
                <tr>
                  <th>項目</th>
                  <th v-for="key in queryMetricsColKeys" :key="key">{{key}}</th>
                </tr>
                </thead>
                <tbody>
                <tr v-for="key in queryMetricsRowKeys" :key="key">
                  <td>{{key}}</td>
                  <td v-for="colKey in queryMetricsColKeys">{{queryClass.metrics[key][colKey]}}</td>
                </tr>
                </tbody>
              </table>
            </div>
          </details>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.container {
  padding: 16px;
}

table {
  margin: 8px 0;
  border-collapse: collapse;
}

td, th {
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

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  data() {
    return {
      summary: "Loading ...",
      slowData: undefined,
      sort: 'total_time',
      queryMetricsRowKeys: {
        Lock_time: 'Lock_time',
        Query_length: 'Query_length',
        Query_time: 'Query_time',
        Rows_examined: 'Rows_examined',
        Rows_sent: 'Rows_sent'
      },
      queryMetricsColKeys: {
        sum: 'sum',
        avg: 'avg',
        max: 'max',
        median: 'median',
        min: 'min',
        pct: 'pct',
        pct_95: 'pct_95',
        stddev: 'stddev'
      }
    };
  },
  computed: {
    classes() {
      if (!this.slowData) return []
      const arr = [...this.slowData.classes]
      switch(this.sort) {
        case 'count':
          return arr.sort((a, b) => b.query_count - a.query_count)
        case 'total_time':
          return arr.sort((a, b) => b.metrics['Query_time']['sum'] - a.metrics['Query_time']['sum'])
        case 'avg_time':
          return arr.sort((a, b) => b.metrics['Query_time']['avg'] - a.metrics['Query_time']['avg'])
        case 'rows_s/e':
          return arr.sort((a, b) => (a.metrics['Rows_sent']['sum']/a.metrics['Rows_examined']['sum']) - (b.metrics['Rows_sent']['sum']/b.metrics['Rows_examined']['sum']))
        default:
          return arr
      }
    }
  },
  methods: {
    copy(text: string) {
      window.navigator.clipboard.writeText(text)
    }
  },
  async beforeCreate() {
    const resp = await fetch(`/api/slowlog/logs/${this.$route.params.id}`);
    this.slowData = await resp.json();
  },
});
</script>

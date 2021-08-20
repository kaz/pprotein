<template>
  <section>
    <pre v-if="!logData">Loading ...</pre>
    <div v-else class="container">
      <div>ヘッダーの項目をクリックでソートできます</div>
      <br />
      <table border="1">
        <thead>
          <tr>
            <th v-for="h in header" :key="h" @click="sortBy(h)">
              <span
                >{{ h
                }}<template v-if="sort == h"
                  ><template v-if="order == 'asc'">↑</template>
                  <template v-if="order == 'desc'">↓</template>
                </template></span
              >
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="d in sortedLogData" :key="d.toString()">
            <td v-for="(v, idx) in d" :key="header[idx]">{{ v }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import Pase, { ParseResult } from "papaparse";

const isNumeric = (v: any) => {
  return !isNaN(+v);
};

export default defineComponent({
  data(): {
    logData: ParseResult<string[]>;
    sort: string;
    order: "desc" | "asc";
  } {
    return {
      logData: {} as Pase.ParseResult<string[]>,
      sort: "count",
      order: "desc",
    };
  },
  computed: {
    header: function () {
      if (!this.logData.data) {
        return [];
      }
      let data: string[][];
      data = this.logData.data;
      return data[0];
    },
    sortedLogData: function () {
      if (!this.logData.data) {
        return [];
      }
      let data: string[][];
      data = this.logData.data.slice(1);
      let index = this.header.findIndex((h: string) => h === this.sort);
      if (index == -1) {
        index = 0;
      }

      data.sort((a: string[], b: string[]): number => {
        if (isNumeric(a[index])) {
          return +a[index] - +b[index];
        } else {
          if (a[index] < b[index]) {
            return -1;
          } else if (a[index] > b[index]) {
            return 1;
          } else {
            return 0;
          }
        }
      });

      if (this.order == "desc") {
        return data.reverse();
      } else {
        return data;
      }
    },
  },
  async beforeCreate() {
    const resp = await fetch(`/api/httplog/${this.$route.params.id}`);
    const text = await resp.text();
    const logData = Pase.parse<string[]>(text, { skipEmptyLines: true });
    this.logData = logData;
  },
  methods: {
    sortBy(key: string) {
      if (this.sort === key) {
        this.order = this.order == "desc" ? "asc" : "desc";
      } else {
        this.sort = key;
        this.order = "desc";
      }
    },
  },
});
</script>

<style scoped lang="scss">
pre {
  margin: 0;
  padding: 2em;
  overflow: auto;
  flex: 1 0 auto;
}

.container {
  overflow: scroll;
}

table {
  border-collapse: collapse;
}

td,
th {
  padding: 4px;
}

th > span {
  cursor: pointer;
}
</style>

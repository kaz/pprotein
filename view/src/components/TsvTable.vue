<template>
  <section>
    <table border="1">
      <thead>
        <tr>
          <th v-for="(h, i) in header" :key="h" @click="sortBy(i)">
            {{ h }}
            <span v-if="sortColumn === i" class="sortOrder">
              {{ sortOrder == "desc" ? "▼" : "▲" }}
            </span>
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="d in sortedData" :key="d.toString()">
          <td v-for="(v, i) in d" :key="header[i]">{{ v }}</td>
        </tr>
      </tbody>
    </table>
  </section>
</template>

<script lang="ts">
import { parse } from "papaparse";
import { defineComponent } from "vue";

const isNumeric = (v: string) => {
  return !isNaN(+v);
};

export default defineComponent({
  props: {
    tsv: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      sortColumn: 0,
      sortOrder: "desc",
    };
  },
  computed: {
    rows() {
      return parse<string[]>(this.tsv, { skipEmptyLines: true }).data;
    },
    header() {
      return this.rows[0] || [];
    },
    sortedData() {
      const data = this.rows.slice(1).sort((as, bs) => {
        const [a, b] = [as[this.sortColumn], bs[this.sortColumn]];
        if (isNumeric(a)) {
          return +a - +b;
        } else {
          if (a < b) {
            return -1;
          } else if (a > b) {
            return 1;
          } else {
            return 0;
          }
        }
      });

      if (this.sortOrder === "desc") {
        return data.reverse();
      } else {
        return data;
      }
    },
  },
  methods: {
    sortBy(column: number) {
      if (this.sortColumn === column) {
        this.sortOrder = this.sortOrder === "desc" ? "asc" : "desc";
      } else {
        this.sortColumn = column;
        this.sortOrder = "desc";
      }
    },
  },
});
</script>

<style scoped lang="scss">
section {
  padding: 2em;
  overflow: scroll;
}

table {
  border-collapse: collapse;
}

td,
th {
  padding: 0.5em 1em;
  border: 1px solid #999;
  white-space: nowrap;
}

th {
  cursor: pointer;

  .sortOrder {
    font-size: 0.4em;
    color: gray;
  }
}
</style>

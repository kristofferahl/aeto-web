<script setup>
import { RouterLink } from 'vue-router'
</script>

<script>
const walk = (obj, prefix, result) => {
  for (let k in obj) {
    if (k === 'managedFields') {
      continue
    }
    if (typeof obj[k] == 'object' && obj[k] !== null && !Array.isArray(obj[k])) {
      result = walk(obj[k], prefix + '.' + k, result)
    } else {
      result[prefix + '.' + k] = obj[k]
    }
  }
  return result
}

export default {
  props: {
    resourceType: String,
    namespace: String,
    name: String,
    fields: Array,
    filter: Function
  },

  data() {
    return {
      error: null,
      resource: null,
      items: null
    }
  },
  methods: {
    async fetchData(namespace, name) {
      try {
        let getOne = namespace && name
        let url = `/api/${this.resourceType.toLowerCase()}`
        if (getOne) {
          url += `/${namespace}/${name}`
        }
        const response = await fetch(url)
        const data = await response.json()
        if (getOne) {
          this.items = null
          this.resource = data
          console.log('details', this.resource)
        } else {
          this.resource = null
          this.items = data.items.filter((i) => (this.filter ? this.filter(i) : i))
          console.log('list', this.items)
        }
      } catch (e) {
        this.error = e
      }
    },
    walkResource(r, mode) {
      return this.fields
        .filter((f) => !f.mode || mode === f.mode)
        .map((f) => {
          return {
            label: f.label,
            value: f.format(r),
            linkTo: f.linkTo ? f.linkTo + 's' : null
          }
        })
    }
  },

  created() {
    this.$watch(
      () => this.$route.params,
      (toParams) => {
        this.fetchData(toParams.namespace, toParams.name)
      }
    )
  },

  mounted() {
    this.fetchData(this.namespace, this.name)
  }
}
</script>

<template>
  <h2>{{ resourceType }} {{ resource ? `: ${namespace}/${name}` : null }}</h2>
  <div v-if="!resource && !items && !error" class="loading">Loading...</div>
  <div v-if="error" class="error">{{ error }}</div>
  <div v-if="resource">
    <table>
      <tbody>
        <tr v-for="f in walkResource(resource, 'details')">
          <td>{{ f.label }}</td>
          <td>{{ f.value }}</td>
        </tr>
      </tbody>
    </table>
  </div>
  <div v-if="items">
    <table>
      <thead>
        <tr>
          <th>Resource</th>
          <th v-for="f in walkResource(items[0], 'list')">
            {{ f.label }}
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="i in items">
          <td>
            <RouterLink
              :to="`/${resourceType.toLowerCase()}/${i.metadata.namespace}/${i.metadata.name}`"
              >{{ i.metadata.namespace + '/' + i.metadata.name }}</RouterLink
            >
          </td>
          <td v-for="f in walkResource(i, 'list')">
            <RouterLink v-if="f.linkTo" :to="`/${f.linkTo.toLowerCase()}/${f.value}`">{{
              f.value
            }}</RouterLink>
            <span v-else>{{ f.value }}</span>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

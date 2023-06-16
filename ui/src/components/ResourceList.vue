<script setup>
import { RouterLink } from 'vue-router'
import { parseISO, formatDistance } from 'date-fns'
</script>

<script>
const capitalize = (s) => {
  return s.charAt(0).toUpperCase() + s.slice(1)
}

export default {
  props: {
    resourceType: String,
    namespace: String,
    name: String,
    fields: Array,
    status: Array
  },

  data() {
    return {
      error: null,
      resource: null,
      items: null,
      viewRaw: false
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
          this.items = data.items
          console.log('list', this.items)
        }
      } catch (e) {
        this.error = e
      }
    },
    walkResource(r, mode) {
      let props = Object.keys(r).map((k) => {
        let override = (this.fields && this.fields.length ? this.fields : []).find(
          (f) => f.key === k && (!f.mode || mode === f.mode)
        )

        if (override && override.hide === true) {
          return null
        }

        if (override?.select) {
          return override.select.map((c) => {
            return {
              ...c,
              key: `${k}.${c.key}`,
              label: c.label ?? capitalize(c.key)
            }
          })
        }

        let def = {
          key: k,
          label: capitalize(k),
          format: (v) => v[k]
        }
        return [{ ...def, ...override }]
      })

      return props
        .flat()
        .filter((f) => f !== null && (!f.mode || mode === f.mode))
        .map((f) => {
          return {
            key: f.key,
            label: f.label,
            value: r != null ? f.format(r) : null,
            linkTo: f.linkTo ? f.linkTo + 's' : null,
            pre: f.pre,
            code: f.code
          }
        })
    },
    walkStatus(s, mode) {
      let props = Object.keys(s)
        .filter((k) => k !== 'conditions')
        .map((k) => {
          let override = (this.status && this.status.length ? this.status : []).find(
            (f) => f.key === k && (!f.mode || mode === f.mode)
          )
          if (override && override.hide === true) {
            return null
          }
          let def = {
            label: capitalize(k),
            format: () => s[k]
          }
          return { ...def, ...override }
        })

      return props
        .filter((f) => f !== null && (!f.mode || mode === f.mode))
        .map((f) => {
          return {
            label: f.label,
            value: s != null ? f.format(s) : null,
            linkTo: f.linkTo ? f.linkTo + 's' : null,
            pre: f.pre,
            code: f.code
          }
        })
    },
    statusConditionHeaders(items) {
      return [
        ...new Set(
          items
            .filter((i) => i.status?.conditions)
            .map((i) => i.status.conditions.map((c) => c.type))
            .flat()
        )
      ]
    },
    walkStatusConditions(s, format) {
      return s.conditions
        ? s.conditions
            .map((c) => (format === true ? this.formatConditions(c) : c))
            .filter((kv) => kv !== null)
        : []
    },
    formatConditions(c) {
      return {
        label: c.type,
        value: Object.keys(c)
          .filter((k) => k !== 'type' && c[k] !== '')
          .map((k) => `${k}=${c[k]}`)
          .join('<br />')
      }
    },
    walkMetadata(m) {
      return Object.keys(m)
        .map((k) => this.formatMetadata(k, m[k]))
        .filter((kv) => kv !== null)
    },
    formatMetadata(k, v) {
      switch (k) {
        case 'creationTimestamp':
          return {
            label: capitalize(k),
            value: formatDistance(parseISO(v), new Date(), { addSuffix: true }) + ` ( ${v} )`
          }
        case 'annotations':
          return {
            label: capitalize(k),
            value: Object.keys(v)
              .filter((a) => a !== 'kubectl.kubernetes.io/last-applied-configuration')
              .map((a) => `${a}=${v[a]}`)
              .join('<br />')
          }
        case 'labels':
          return {
            label: capitalize(k),
            value: Object.keys(v)
              .map((l) => `${l}=${v[l]}`)
              .join('<br />')
          }
        case 'finalizers':
          return { label: capitalize(k), value: v.join('<br />') }
        case 'managedFields':
          return null
        default:
          return { label: capitalize(k), value: v }
      }
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
  <div v-if="resource" class="details row">
    <section class="column">
      <button @click="viewRaw = !viewRaw" class="button button-small">
        {{ viewRaw ? 'View normal' : 'View raw' }}
      </button>
    </section>
  </div>
  <div v-if="resource && viewRaw" class="details row">
    <section class="column">
      <pre>{{ resource }}</pre>
    </section>
  </div>
  <div v-if="resource && !viewRaw" class="details row">
    <section class="column column-75">
      <table>
        <tbody>
          <tr>
            <th colspan="2">Resource</th>
          </tr>
          <tr v-for="f in walkResource(resource.spec, 'details')" :key="f.key">
            <td :title="f.key">{{ f.label }}</td>
            <td>
              <RouterLink v-if="f.linkTo" :to="`/${f.linkTo.toLowerCase()}/${f.value}`">{{
                f.value
              }}</RouterLink>
              <pre v-else-if="f.pre">{{ f.value }}</pre>
              <pre v-else-if="f.code" class="withcode"><code>{{ f.value }}</code></pre>
              <span v-else>{{ f.value }}</span>
            </td>
          </tr>
          <tr>
            <th colspan="2">Status</th>
          </tr>
          <div v-if="resource.status">
            <tr v-for="f in walkStatus(resource.status, 'details')" :key="f.label">
              <td>{{ f.label }}</td>
              <td>
                <RouterLink v-if="f.linkTo" :to="`/${f.linkTo.toLowerCase()}/${f.value}`">{{
                  f.value
                }}</RouterLink>
                <pre v-else-if="f.pre">{{ f.value }}</pre>
                <pre v-else-if="f.code" class="withcode"><code>{{ f.value }}</code></pre>
                <span v-else>{{ f.value }}</span>
              </td>
            </tr>
          </div>
          <tr>
            <th colspan="2">Conditions</th>
          </tr>
          <tr v-for="kv in walkStatusConditions(resource.status, true)" :key="kv.label">
            <td>{{ kv.label }}</td>
            <td>
              <pre v-html="kv.value"></pre>
            </td>
          </tr>
        </tbody>
      </table>
    </section>

    <section class="column column-25">
      <table>
        <tbody>
          <tr>
            <th colspan="2">Metadata</th>
          </tr>
          <tr v-for="kv in walkMetadata(resource.metadata)" :key="kv.label">
            <td>{{ kv.label }}</td>
            <td><span v-html="kv.value"></span></td>
          </tr>
        </tbody>
      </table>
    </section>
  </div>
  <div v-if="items && items.length">
    <table>
      <thead>
        <tr>
          <th>Resource</th>
          <th v-for="f in walkResource(items[0].spec, 'list')" :key="f.key" :title="f.key">
            {{ f.label }}
          </th>
          <th v-if="items.some((i) => i.status?.status !== undefined)">Status</th>
          <th v-for="(h, i) in statusConditionHeaders(items)" :key="i">
            {{ h }}
          </th>
          <th>Created</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="i in items" :key="i.metadata.uid">
          <td>
            <RouterLink
              :to="`/${resourceType.toLowerCase()}/${i.metadata.namespace}/${i.metadata.name}`"
              >{{ i.metadata.namespace + '/' + i.metadata.name }}</RouterLink
            >
          </td>
          <td v-for="(f, i) in walkResource(i.spec, 'list')" :key="i">
            <RouterLink v-if="f.linkTo" :to="`/${f.linkTo.toLowerCase()}/${f.value}`">{{
              f.value
            }}</RouterLink>
            <span v-else>{{ f.value }}</span>
          </td>
          <td v-if="i.status?.status !== undefined">
            {{ i.status.status }}
          </td>
          <td v-for="kv in walkStatusConditions(i.status)" :title="kv.key" :key="kv.key">
            {{ kv.status }}
          </td>
          <td>
            {{
              formatDistance(parseISO(i.metadata.creationTimestamp), new Date(), {
                addSuffix: true
              })
            }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

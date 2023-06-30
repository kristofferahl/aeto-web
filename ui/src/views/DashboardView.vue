<script setup>
import { parseISO, formatDistance } from 'date-fns'
</script>

<script>
export default {
  data() {
    return {
      error: null,
      dashboard: {},
      source: null
    }
  },
  methods: {
    async fetchData(namespace, name) {
      try {
        let url = `/api/dashboard`
        const response = await fetch(url)
        const data = await response.json()
        data.changes = data.changes.reverse()
        this.dashboard = data
      } catch (e) {
        this.error = e
      }
    }
  },

  mounted() {
    this.fetchData()
    this.source = new EventSource('/api/sse')
    this.source.onmessage = (e) => {
      console.log('Server sent event:')
      console.log(e)
    }
    this.source.onerror = (err) => {
      console.error(err)
      this.source.close()
    }
  }
}
</script>

<template>
  <h2>Dashboard</h2>
  <div class="dashboard">
    <div class="row">
      <div class="column column-20">
        <div class="card stats">
          <div class="stats-content">{{ dashboard.tenants }}</div>
          <div class="stats-title">Tenants</div>
        </div>
      </div>
      <div class="column column-40">
        <div class="card">
          <h3>Resource Changes</h3>
          <ul>
            <li v-for="c in dashboard.changes">
              {{ c.change }} {{ c.type }} {{ c.resource }} ({{
                formatDistance(parseISO(c.ts), new Date(), { addSuffix: true })
              }})
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

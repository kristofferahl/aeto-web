<script setup>
import { parseISO, formatDistance } from 'date-fns'
</script>

<script>
export default {
  data() {
    return {
      error: null,
      dashboard: {}
    }
  },
  methods: {
    async fetchData(namespace, name) {
      try {
        let url = `/api/dashboard`
        const response = await fetch(url)
        const data = await response.json()
        this.dashboard = data
      } catch (e) {
        this.error = e
      }
    }
  },

  mounted() {
    this.fetchData()
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
    </div>
  </div>
</template>

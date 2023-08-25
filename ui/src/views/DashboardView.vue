<script setup>
import { parseISO, formatDistance } from 'date-fns'
import { useEventStreamStore } from '../stores/eventStream'

const { eventStream } = useEventStreamStore()
</script>

<script>
function uuidv4() {
  return ([1e7] + -1e3 + -4e3 + -8e3 + -1e11).replace(/[018]/g, (c) =>
    (c ^ (crypto.getRandomValues(new Uint8Array(1))[0] & (15 >> (c / 4)))).toString(16)
  )
}

export default {
  data() {
    return {
      error: null,
      dashboard: {},
      // eventstream: [],
      source: null
    }
  },
  methods: {
    async fetchData() {
      try {
        let url = `/api/dashboard`
        const response = await fetch(url)
        const data = await response.json()
        console.log('dashboard data', data)
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
      <div class="column column-10">
        <div class="card stats">
          <div class="stats-content">{{ dashboard.tenants }}</div>
          <div class="stats-title">Tenants</div>
        </div>
        <div class="card stats">
          <div class="stats-content">{{ dashboard.blueprints }}</div>
          <div class="stats-title">Blueprints</div>
        </div>
        <div class="card stats">
          <div class="stats-content">{{ dashboard.certificates }}</div>
          <div class="stats-title">Certificates</div>
        </div>
        <div class="card stats">
          <div class="stats-content">{{ dashboard.hostedzones }}</div>
          <div class="stats-title">HostedZones</div>
        </div>
        <div class="card stats">
          <div class="stats-content">{{ dashboard.savingspolicies }}</div>
          <div class="stats-title">SavingsPolicies</div>
        </div>
      </div>
      <div class="column column-60">
        <div class="card">
          <h3>Resource Changes</h3>
          <ul>
            <li v-for="e in eventStream.filter((e) =>
              ['ResourceAdded', 'ResourceUpdated', 'ResourceDeleted'].includes(e.type)
            )">
              {{ e.type }}<br />
              {{ e.resource.apiVersion }}/{{ e.resource.kind }}<br />
              {{ e.resource.metadata.namespace }}/{{ e.resource.metadata.name }}<br />
              ({{ formatDistance(parseISO(e.ts), new Date(), { addSuffix: true }) }})
            </li>
          </ul>
        </div>
      </div>
      <div class="column column-30">
        <div class="card">
          <h3>Kubernetes Events</h3>
          <ul>
            <li v-for="e in eventStream.filter((e) =>
              [
                'KubernetesEventAdded',
                'KubernetesEventUpdated',
                'KubernetesEventDeleted'
              ].includes(e.type)
            )">
              {{ e.message }}<br />
              {{ e.resource.apiVersion }}/{{ e.resource.kind }}<br />
              {{ e.resource.namespace }}/{{ e.resource.name }}<br />
              ({{ formatDistance(parseISO(e.ts), new Date(), { addSuffix: true }) }})
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

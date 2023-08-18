<script setup>
import { parseISO, formatDistance } from 'date-fns'
</script>

<script>
function uuidv4() {
  return ([1e7] + -1e3 + -4e3 + -8e3 + -1e11).replace(/[018]/g, (c) =>
    (c ^ (crypto.getRandomValues(new Uint8Array(1))[0] & (15 >> (c / 4)))).toString(16)
  )
}

function prepend(a, v) {
  var na = a.slice()
  na.unshift(v)
  return na
}

export default {
  data() {
    return {
      error: null,
      dashboard: {},
      eventstream: [],
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
    this.source = new EventSource('/api/sse?cid=' + uuidv4().substring(0, 6))
    console.log('Connecting to server')
    this.source.onmessage = (e) => {
      console.log('New event', e)
      this.eventstream = prepend(this.eventstream, JSON.parse(e.data))
    }
    this.source.onerror = (err) => {
      console.error('Error', err)
      this.source.close()
    }
  },
  beforeUnmount() {
    console.log('Disconnecting from server')
    this.source.close()
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
      <div class="column column-40">
        <div class="card">
          <h3>Resource Changes</h3>
          <ul>
            <li
              v-for="e in eventstream.filter((e) =>
                ['ResourceAdded', 'ResourceUpdated', 'ResourceDeleted'].includes(e.type)
              )"
            >
              {{ e.type }}<br />
              {{ e.resource.apiVersion }}/{{ e.resource.kind }}<br />
              {{ e.resource.metadata.namespace }}/{{ e.resource.metadata.name }}<br />
              ({{ formatDistance(parseISO(e.ts), new Date(), { addSuffix: true }) }})
            </li>
          </ul>
        </div>
      </div>
      <div class="column column-40">
        <div class="card">
          <h3>Kubernetes Events</h3>
          <ul>
            <li
              v-for="e in eventstream.filter((e) =>
                [
                  'KubernetesEventAdded',
                  'KubernetesEventUpdated',
                  'KubernetesEventDeleted'
                ].includes(e.type)
              )"
            >
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

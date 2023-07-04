<script setup>
import { parseISO, formatDistance } from 'date-fns'
</script>

<script>
function uuidv4() {
  return ([1e7]+-1e3+-4e3+-8e3+-1e11).replace(/[018]/g, c =>
    (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
  );
}

function prepend(a, v) {
  var na = a.slice();
  na.unshift(v);
  return na;
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
    this.source = new EventSource('/api/sse?cid='+uuidv4().substring(0,6))
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
      </div>
      <div class="column column-40">
        <div class="card">
          <h3>Resource Changes</h3>
          <ul>
            <li v-for="c in eventstream">
              <span v-if="['ResourceAdded','ResourceUpdated','ResourceDeleted'].includes(c.type)">
                {{ c.type }}<br />
                {{ c.resource.apiVersion }}/{{ c.resource.kind }}<br />
                {{ c.resource.metadata.namespace }}/{{ c.resource.metadata.name }}<br />
                ({{
                  formatDistance(parseISO(c.ts), new Date(), { addSuffix: true })
                }})
              </span>
              <span v-else>
                {{ c.type }} {{ c.payload }} ({{
                  formatDistance(parseISO(c.ts), new Date(), { addSuffix: true })
                }})
              </span>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      error: null,
      tenants: null
    }
  },
  methods: {
    async fetchData() {
      try {
        const response = await fetch('/api/tenants')
        const data = await response.json()
        this.tenants = data.items
        console.log('tenants', this.tenants)
      } catch (e) {
        this.error = e
      }
    }
  },

  mounted() {
    console.log('initially loading data')
    this.fetchData()
  }
}
</script>

<template>
  <h2>Tenants</h2>
  <div v-if="!tenants && !error" class="loading">Loading...</div>
  <div v-if="error" class="error">{{ error }}</div>
  <div v-if="tenants" class="content">
    <table>
      <thead>
        <th>Name</th>
        <th>Fullname</th>
        <th>Blueprint</th>
        <th>Namespace</th>
        <th>Status</th>
        <th>ResourceSet</th>
        <th>Events</th>
      </thead>
      <tbody>
        <tr v-for="t in tenants">
          <td>{{ t.metadata.name }}</td>
          <td>{{ t.spec.name }}</td>
          <td>{{ t.status.blueprint }}</td>
          <td>{{ t.status.namespace }}</td>
          <td>{{ t.status.status }}</td>
          <td>{{ t.status.resourceSet }}</td>
          <td>{{ t.status.events }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
table thead th {
  text-align: left;
  font-weight: bold;
}

table thead th,
table tbody td {
  padding: 0.2em 1em 0.2em 0;
}
</style>

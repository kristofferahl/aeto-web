<script setup>
import Resources from '../components/Resources.vue'
</script>

<script>
export default {
  computed: {
    namespace() {
      return this.$route.params.namespace
    },
    name() {
      return this.$route.params.name
    },
    fields() {
      return [
        { key: 'ts', label: 'Timestamp' },
        { key: 'events', mode: 'list', format: (spec) => spec.events.length },
        {
          key: 'events',
          mode: 'details',
          format: (spec) => spec.events.map((e) => JSON.parse(e.raw)),
          pre: true
        }
      ]
    },
    status() {
      return []
    }
  }
}
</script>

<template>
  <Resources
    resourceType="EventStreamChunks"
    :namespace="namespace"
    :name="name"
    :fields="fields"
    :status="status"
  />
</template>

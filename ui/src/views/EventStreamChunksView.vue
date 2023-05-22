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
        { label: 'Name', format: (r) => r.metadata.name, mode: 'details' },
        { label: 'Id', format: (r) => r.spec.id },
        { label: 'Version', format: (r) => r.spec.version },
        { label: 'Timestamp', format: (r) => r.spec.ts },
        { label: 'Events', format: (r) => r.spec.events.length },
      ]
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
  />
</template>

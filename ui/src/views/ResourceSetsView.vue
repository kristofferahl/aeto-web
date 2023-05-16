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
        { label: 'Namespace', format: (r) => r.metadata.namespace, mode: 'details' },
        { label: 'Created', format: (r) => r.metadata.creationTimestamp, mode: 'details' },
        {
          label: 'Active',
          format: (r) => r.status.conditions.find((c) => c.type === 'Active').status
        },
        {
          label: 'Ready',
          format: (r) => r.status.conditions.find((c) => c.type === 'Ready').message
        },
        { label: 'Resources', format: (r) => r.spec.resources.length },
        { label: 'Status', format: (r) => r.status.status }
      ]
    }
  }
}
</script>

<template>
  <Resources resourceType="ResourceSets" :namespace="namespace" :name="name" :fields="fields" />
</template>

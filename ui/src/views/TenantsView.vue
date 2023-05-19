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
        { label: 'Fullname', format: (r) => r.spec.name },
        { label: 'Blueprint', format: (r) => r.status.blueprint, linkTo: 'blueprint' },
        { label: 'Namespace', format: (r) => r.status.namespace },
        { label: 'ResourceSet', format: (r) => r.status.resourceSet, linkTo: 'resourceset' },
        { label: 'Events', format: (r) => r.status.events },
        { label: 'Created', format: (r) => r.metadata.creationTimestamp, mode: 'details' },
        {
          label: 'Ready',
          format: (r) => r.status.conditions.find((c) => c.type === 'Ready').status,
          mode: 'list'
        },
        {
          label: 'Ready',
          format: (r) => r.status.conditions.find((c) => c.type === 'Ready').message,
          mode: 'details'
        },
        { label: 'Status', format: (r) => r.status.status }
      ]
    }
  }
}
</script>

<template>
  <Resources resourceType="Tenants" :namespace="namespace" :name="name" :fields="fields" />
</template>

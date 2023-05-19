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
        { label: 'ResourceNamePrefix', format: (r) => r.spec.resourceNamePrefix },
        { label: 'Resources', format: (r) => r.spec.resources.length },
        { label: 'Created', format: (r) => r.metadata.creationTimestamp, mode: 'details' }
      ]
    }
  }
}
</script>

<template>
  <Resources resourceType="Blueprints" :namespace="namespace" :name="name" :fields="fields" />
</template>

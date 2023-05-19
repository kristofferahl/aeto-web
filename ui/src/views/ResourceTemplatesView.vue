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
        { label: 'Parameters', format: (r) => r.spec.parameters, mode: 'details' },
        { label: 'Parameters', format: (r) => r.spec.parameters.length, mode: 'list' },
        { label: 'Resources', format: (r) => r.spec.resources.length },
        { label: 'Rules', format: (r) => r.spec.rules }
      ]
    }
  }
}
</script>

<template>
  <Resources resourceType="ResourceTemplates" :namespace="namespace" :name="name" :fields="fields" />
</template>

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
      return []
    },
    status() {
      return [
        { key: 'blueprint', format: (status) => status.blueprint, linkTo: 'blueprint' },
        { key: 'resourceSet', format: (status) => status.resourceSet, linkTo: 'resourceset' }
      ]
    }
  }
}
</script>

<template>
  <Resources
    resourceType="Tenants"
    :namespace="namespace"
    :name="name"
    :fields="fields"
    :status="status"
  />
</template>

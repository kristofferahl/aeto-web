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
      return [{ key: 'suspended', format: (spec) => spec.suspended.map((s) => s).join(',') }]
    },
    status() {
      return []
    }
  }
}
</script>

<template>
  <Resources
    resourceType="SavingsPolicies"
    :namespace="namespace"
    :name="name"
    :fields="fields"
    :status="status"
  />
</template>

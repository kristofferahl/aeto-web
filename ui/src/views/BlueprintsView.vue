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
        { key: 'resources', mode: 'list', format: (spec) => spec.resources.length },
        { key: 'resources', mode: 'details', format: (spec) => spec.resources, pre: true }
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
    resourceType="Blueprints"
    :namespace="namespace"
    :name="name"
    :fields="fields"
    :status="status"
  />
</template>

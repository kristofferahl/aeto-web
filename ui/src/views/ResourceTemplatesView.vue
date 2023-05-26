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
        { key: 'parameters', mode: 'details', format: (spec) => spec.parameters, pre: true },
        { key: 'parameters', mode: 'list', format: (spec) => spec.parameters.length },
        { key: 'resources', mode: 'list', format: (spec) => spec.resources.length },
        { key: 'resources', mode: 'details', format: (spec) => spec.resources, pre: true },
        { key: 'rules', format: (spec) => spec.rules, pre: true }
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
    resourceType="ResourceTemplates"
    :namespace="namespace"
    :name="name"
    :fields="fields"
    :status="status"
  />
</template>

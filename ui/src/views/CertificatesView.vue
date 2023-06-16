<script setup>
import Resources from '../components/ResourceList.vue'
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
        { key: 'validation', format: (spec) => Object.keys(spec.validation)[0], mode: 'list' },
        { key: 'validation', format: (spec) => spec.validation, mode: 'details', pre: true },
        { key: 'tags', mode: 'details', format: (spec) => spec.tags, pre: true },
        { key: 'tags', mode: 'list', format: (spec) => Object.keys(spec.tags).length }
      ]
    },
    status() {
      return [{ key: 'arn', format: (status) => status.arn, mode: 'details', code: true }]
    }
  }
}
</script>

<template>
  <Resources
    resourceType="Certificates"
    :namespace="namespace"
    :name="name"
    :fields="fields"
    :status="status"
  />
</template>

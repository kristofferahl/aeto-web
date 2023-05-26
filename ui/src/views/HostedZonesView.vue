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
        {
          key: 'connectWith',
          select: [
            { key: 'name', label: 'Connect with', format: (spec) => spec.connectWith.name },
            { key: 'ttl', label: 'TTL', format: (spec) => spec.connectWith.ttl }
          ]
        },
        { key: 'tags', mode: 'details', format: (spec) => spec.tags, pre: true },
        { key: 'tags', mode: 'list', hide: true }
      ]
    },
    status() {
      return [{ key: 'id', format: (status) => status.id, code: true }]
    }
  }
}
</script>

<template>
  <Resources
    resourceType="HostedZones"
    :namespace="namespace"
    :name="name"
    :fields="fields"
    :status="status"
  />
</template>

<script setup>
import { setBlockTracking } from 'vue'
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
        {
          key: 'certificates',
          mode: 'list',
          hide: true
        },
        {
          key: 'certificates',
          label: 'Certificate Selector',
          mode: 'details',
          format: (spec) => spec.certificates.selector,
          pre: true
        },
        {
          key: 'ingress',
          mode: 'list',
          label: 'Ingress Connector',
          format: (spec) => spec.ingress.connector
        },
        {
          key: 'ingress',
          mode: 'details',
          select: [
            {
              key: 'connector',
              label: 'Ingress Connector',
              format: (spec) => spec.ingress.connector
            },
            {
              key: 'selector',
              label: 'Ingress Selector',
              format: (spec) => spec.ingress.selector,
              pre: true
            }
          ]
        }
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
    resourceType="CertificateConnectors"
    :namespace="namespace"
    :name="name"
    :fields="fields"
    :status="status"
  />
</template>

import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useEventStreamStore = defineStore('eventStream', () => {
  const eventStream = ref([])

  function prependEvent(v) {
    var na = eventStream.value.slice()
    na.unshift(v)
    eventStream.value = na
  }

  return { eventStream, prependEvent }
})

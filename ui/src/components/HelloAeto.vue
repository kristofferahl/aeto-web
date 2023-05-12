<script setup>
defineProps({
  msg: {
    type: String,
    required: true
  }
})
</script>

<template>
  <div class="greetings">
    <h1 class="green">aeto</h1>
    <h2>aws-eks-tenant-operator</h2>
    <h3>{{ msg }}</h3>
    <p>
      Find out more on
      <a href="https://github.com/kristofferahl/aeto" target="_blank" rel="noopener">github</a>
    </p>
  </div>
</template>

<style scoped>
h1 {
  font-weight: 500;
  font-size: 2.6rem;
  top: -10px;
}

h2 {
  top: -20px;
}

h3 {
  font-size: 1.2rem;
}

.greetings h1,
.greetings h2,
.greetings h3,
.greetings p {
  text-align: center;
}

@media (min-width: 1024px) {
  .greetings h1,
  .greetings h2,
  .greetings h3,
  .greetings p {
    text-align: left;
  }
}
</style>

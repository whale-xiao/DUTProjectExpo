<script setup>
defineProps({
  block: { type: Object, required: true },
})
</script>

<template>
  <div v-if="block.type === 'heading'" class="block block-heading">
    {{ block.text }}
  </div>
  <p v-else-if="block.type === 'paragraph'" class="block block-paragraph">
    {{ block.text }}
  </p>
  <figure v-else-if="block.type === 'image'" class="block block-image">
    <img v-if="block.url" :src="block.url" :alt="block.caption || '项目配图'" loading="lazy" decoding="async" />
    <div v-else class="img-ph">图片占位</div>
    <figcaption v-if="block.caption">{{ block.caption }}</figcaption>
  </figure>
</template>

<style scoped>
.block {
  color: var(--ink);
}
.block-heading {
  font-family: var(--serif);
  font-size: 22px;
  font-weight: 700;
  margin: 32px 0 0;
}
.block-paragraph {
  font-family: var(--serif);
  font-size: 17px;
  line-height: 1.9;
  letter-spacing: 0.02em;
  margin: 16px 0 0;
}
.block-image {
  margin: 24px 0 0;
}
.block-image img {
  max-width: 100%;
  max-height: 540px;
  border-radius: var(--r-md);
  display: block;
  margin: 0 auto;
}
.block-image figcaption {
  text-align: center;
  color: var(--muted);
  font-size: 13px;
  margin-top: 8px;
}
.img-ph {
  border: 1px dashed var(--hairline);
  border-radius: var(--r-md);
  padding: 56px 0;
  text-align: center;
  color: var(--muted);
  background: var(--bg);
}
</style>

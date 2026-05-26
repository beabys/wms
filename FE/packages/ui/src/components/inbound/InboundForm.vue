<script setup lang="ts">
import { ref, watch } from 'vue'
defineOptions({ name: 'InboundForm' })

export interface InboundProduct {
  id: string
  sku: string
  name: string
  quantity: number
  length?: number
  width?: number
  height?: number
  weight?: number
  packaging: 'box' | 'pallet' | 'envelope' | 'loose'
}

const props = defineProps<{
  products?: InboundProduct[]
  errors?: Record<string, string>
  loading?: boolean
}>()

const emit = defineEmits<{
  addProduct: []
  removeProduct: [id: string]
  updateProduct: [id: string, field: string, value: unknown]
  submit: []
}>()

const localProducts = ref<InboundProduct[]>(props.products || [])

watch(() => props.products, (val) => {
  if (val) localProducts.value = val
}, { immediate: true })
</script>

<template>
  <div class="bg-white border rounded-lg p-6 space-y-6 dark:bg-gray-800 dark:border-gray-700">
    <h2 class="text-xl font-semibold text-gray-800 dark:text-gray-100">New Inbound Shipment</h2>

    <div v-if="errors && errors._general" class="text-sm text-red-600 bg-red-50 border border-red-200 rounded px-3 py-2 dark:text-red-400 dark:bg-red-900/30 dark:border-red-800">
      {{ errors._general }}
    </div>

    <div
      v-for="(product, idx) in localProducts"
      :key="product.id"
      class="p-4 border rounded-lg space-y-3 dark:border-gray-600"
      :class="{ 'border-red-300 bg-red-50 dark:border-red-800 dark:bg-red-900/20': errors && errors[product.id] }"
    >
      <div class="flex items-center justify-between">
        <span class="text-sm font-medium text-gray-700 dark:text-gray-300">Product {{ idx + 1 }}</span>
        <button
          class="text-sm text-red-600 hover:text-red-800"
          @click="emit('removeProduct', product.id)"
        >
          Remove
        </button>
      </div>

      <div class="grid grid-cols-2 gap-3">
        <div>
          <label class="block text-xs text-gray-500 mb-1 dark:text-gray-400">SKU</label>
          <input
            :value="product.sku"
            placeholder="SKU-001"
            class="w-full px-2 py-1.5 border rounded text-sm dark:bg-gray-700 dark:text-gray-100 dark:border-gray-600"
            @input="emit('updateProduct', product.id, 'sku', ($event.target as HTMLInputElement).value)"
          />
        </div>
        <div>
          <label class="block text-xs text-gray-500 mb-1 dark:text-gray-400">Name</label>
          <input
            :value="product.name"
            placeholder="Product name"
            class="w-full px-2 py-1.5 border rounded text-sm dark:bg-gray-700 dark:text-gray-100 dark:border-gray-600"
            @input="emit('updateProduct', product.id, 'name', ($event.target as HTMLInputElement).value)"
          />
        </div>
        <div>
          <label class="block text-xs text-gray-500 mb-1 dark:text-gray-400">Quantity</label>
          <input
            :value="product.quantity"
            type="number"
            min="1"
            class="w-full px-2 py-1.5 border rounded text-sm dark:bg-gray-700 dark:text-gray-100 dark:border-gray-600"
            @input="emit('updateProduct', product.id, 'quantity', parseInt(($event.target as HTMLInputElement).value) || 0)"
          />
        </div>
        <div>
          <label class="block text-xs text-gray-500 mb-1 dark:text-gray-400">Packaging</label>
          <select
            :value="product.packaging"
            class="w-full px-2 py-1.5 border rounded text-sm dark:bg-gray-700 dark:text-gray-100 dark:border-gray-600"
            @change="emit('updateProduct', product.id, 'packaging', ($event.target as HTMLSelectElement).value)"
          >
            <option value="box">Box</option>
            <option value="pallet">Pallet</option>
            <option value="envelope">Envelope</option>
            <option value="loose">Loose</option>
          </select>
        </div>
        <div>
          <label class="block text-xs text-gray-500 mb-1 dark:text-gray-400">Length (cm)</label>
          <input
            :value="product.length"
            type="number"
            min="0"
            step="0.1"
            class="w-full px-2 py-1.5 border rounded text-sm dark:bg-gray-700 dark:text-gray-100 dark:border-gray-600"
            @input="emit('updateProduct', product.id, 'length', parseFloat(($event.target as HTMLInputElement).value) || undefined)"
          />
        </div>
        <div>
          <label class="block text-xs text-gray-500 mb-1 dark:text-gray-400">Width (cm)</label>
          <input
            :value="product.width"
            type="number"
            min="0"
            step="0.1"
            class="w-full px-2 py-1.5 border rounded text-sm dark:bg-gray-700 dark:text-gray-100 dark:border-gray-600"
            @input="emit('updateProduct', product.id, 'width', parseFloat(($event.target as HTMLInputElement).value) || undefined)"
          />
        </div>
        <div>
          <label class="block text-xs text-gray-500 mb-1 dark:text-gray-400">Height (cm)</label>
          <input
            :value="product.height"
            type="number"
            min="0"
            step="0.1"
            class="w-full px-2 py-1.5 border rounded text-sm dark:bg-gray-700 dark:text-gray-100 dark:border-gray-600"
            @input="emit('updateProduct', product.id, 'height', parseFloat(($event.target as HTMLInputElement).value) || undefined)"
          />
        </div>
        <div>
          <label class="block text-xs text-gray-500 mb-1 dark:text-gray-400">Weight (kg)</label>
          <input
            :value="product.weight"
            type="number"
            min="0"
            step="0.1"
            class="w-full px-2 py-1.5 border rounded text-sm dark:bg-gray-700 dark:text-gray-100 dark:border-gray-600"
            @input="emit('updateProduct', product.id, 'weight', parseFloat(($event.target as HTMLInputElement).value) || undefined)"
          />
        </div>
      </div>

      <p v-if="errors && errors[product.id]" class="text-xs text-red-500 dark:text-red-400">
        {{ errors[product.id] }}
      </p>
    </div>

    <button
      class="px-4 py-2 border border-dashed border-gray-300 rounded text-sm text-gray-500 hover:border-blue-400 hover:text-blue-600 dark:border-gray-600 dark:text-gray-400 dark:hover:border-blue-400 dark:hover:text-blue-400"
      @click="emit('addProduct')"
    >
      + Add Product
    </button>

    <div class="flex justify-end pt-2">
      <button
        :disabled="loading || localProducts.length === 0"
        class="px-6 py-2 bg-blue-600 text-white rounded text-sm font-medium hover:bg-blue-700 disabled:opacity-60"
        @click="emit('submit')"
      >
        {{ loading ? 'Submitting...' : 'Submit Shipment' }}
      </button>
    </div>
  </div>
</template>

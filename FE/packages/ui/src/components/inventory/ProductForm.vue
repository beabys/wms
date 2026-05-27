<script setup lang="ts">
import { ref, computed } from 'vue'
defineOptions({ name: 'ProductForm' })

export interface ProductFormData {
  sku: string
  name: string
  description: string
  category: string
  unit: string
  weight_kg: number | null
  low_stock_threshold: number | null
}

const props = defineProps<{
  mode?: 'create' | 'edit'
  initial?: Partial<ProductFormData>
  loading?: boolean
  errors?: Record<string, string>
}>()

const emit = defineEmits<{
  submit: [data: ProductFormData]
  cancel: []
}>()

const form = ref<ProductFormData>({
  sku: props.initial?.sku || '',
  name: props.initial?.name || '',
  description: props.initial?.description || '',
  category: props.initial?.category || '',
  unit: props.initial?.unit || 'piece',
  weight_kg: props.initial?.weight_kg ?? null,
  low_stock_threshold: props.initial?.low_stock_threshold ?? 10,
})

const validationErrors = ref<Record<string, string>>({})

const categories = ['', 'Electronics', 'Clothing', 'Food', 'Pharmaceutical', 'Automotive', 'Furniture', 'Other']
const units = ['piece', 'kg', 'g', 'l', 'ml', 'box', 'pallet']

function validate(): boolean {
  const errs: Record<string, string> = {}
  if (!form.value.sku.trim()) errs.sku = 'SKU is required'
  if (!form.value.name.trim()) errs.name = 'Name is required'
  if (form.value.weight_kg !== null && form.value.weight_kg <= 0) errs.weight_kg = 'Weight must be positive'
  validationErrors.value = errs
  return Object.keys(errs).length === 0
}

function handleSubmit(): void {
  if (!validate()) return
  emit('submit', { ...form.value })
}

const isEdit = computed(() => props.mode === 'edit')
</script>

<template>
  <div class="bg-white border rounded-lg p-6 space-y-6 dark:bg-gray-800 dark:border-gray-700">
    <h2 class="text-xl font-semibold text-gray-800 dark:text-gray-100">
      {{ isEdit ? 'Edit Product' : 'Create Product' }}
    </h2>

    <div v-if="errors && errors._general" class="text-sm text-red-600 bg-red-50 border border-red-200 rounded px-3 py-2 dark:text-red-400 dark:bg-red-900/30 dark:border-red-800">
      {{ errors._general }}
    </div>

    <div class="grid grid-cols-2 gap-4">
      <div>
        <label class="block text-xs text-gray-500 mb-1 dark:text-gray-400">SKU *</label>
        <input
          v-model="form.sku"
          :disabled="isEdit"
          placeholder="SKU-001"
          class="w-full px-2 py-1.5 border rounded text-sm dark:bg-gray-700 dark:text-gray-100 dark:border-gray-600"
          :class="{ 'border-red-400': validationErrors.sku }"
        />
        <p v-if="validationErrors.sku" class="text-xs text-red-500 mt-1">{{ validationErrors.sku }}</p>
      </div>
      <div>
        <label class="block text-xs text-gray-500 mb-1 dark:text-gray-400">Name *</label>
        <input
          v-model="form.name"
          placeholder="Product name"
          class="w-full px-2 py-1.5 border rounded text-sm dark:bg-gray-700 dark:text-gray-100 dark:border-gray-600"
          :class="{ 'border-red-400': validationErrors.name }"
        />
        <p v-if="validationErrors.name" class="text-xs text-red-500 mt-1">{{ validationErrors.name }}</p>
      </div>
      <div class="col-span-2">
        <label class="block text-xs text-gray-500 mb-1 dark:text-gray-400">Description</label>
        <textarea
          v-model="form.description"
          placeholder="Product description"
          rows="3"
          class="w-full px-2 py-1.5 border rounded text-sm dark:bg-gray-700 dark:text-gray-100 dark:border-gray-600"
        ></textarea>
      </div>
      <div>
        <label class="block text-xs text-gray-500 mb-1 dark:text-gray-400">Category</label>
        <select
          v-model="form.category"
          class="w-full px-2 py-1.5 border rounded text-sm dark:bg-gray-700 dark:text-gray-100 dark:border-gray-600"
        >
          <option v-for="cat in categories" :key="cat" :value="cat">{{ cat || '-- Select --' }}</option>
        </select>
      </div>
      <div>
        <label class="block text-xs text-gray-500 mb-1 dark:text-gray-400">Unit</label>
        <select
          v-model="form.unit"
          class="w-full px-2 py-1.5 border rounded text-sm dark:bg-gray-700 dark:text-gray-100 dark:border-gray-600"
        >
          <option v-for="u in units" :key="u" :value="u">{{ u }}</option>
        </select>
      </div>
      <div>
        <label class="block text-xs text-gray-500 mb-1 dark:text-gray-400">Weight (kg)</label>
        <input
          v-model.number="form.weight_kg"
          type="number"
          min="0"
          step="0.01"
          placeholder="0.00"
          class="w-full px-2 py-1.5 border rounded text-sm dark:bg-gray-700 dark:text-gray-100 dark:border-gray-600"
          :class="{ 'border-red-400': validationErrors.weight_kg }"
        />
        <p v-if="validationErrors.weight_kg" class="text-xs text-red-500 mt-1">{{ validationErrors.weight_kg }}</p>
      </div>
      <div>
        <label class="block text-xs text-gray-500 mb-1 dark:text-gray-400">Low Stock Threshold</label>
        <input
          v-model.number="form.low_stock_threshold"
          type="number"
          min="0"
          placeholder="10"
          class="w-full px-2 py-1.5 border rounded text-sm dark:bg-gray-700 dark:text-gray-100 dark:border-gray-600"
        />
      </div>
    </div>

    <div class="flex justify-end gap-3 pt-2">
      <button
        class="px-4 py-2 border border-gray-300 rounded text-sm text-gray-600 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-700"
        @click="emit('cancel')"
      >
        Cancel
      </button>
      <button
        :disabled="loading"
        class="px-6 py-2 bg-blue-600 text-white rounded text-sm font-medium hover:bg-blue-700 disabled:opacity-60"
        @click="handleSubmit"
      >
        {{ loading ? 'Saving...' : isEdit ? 'Update Product' : 'Create Product' }}
      </button>
    </div>
  </div>
</template>

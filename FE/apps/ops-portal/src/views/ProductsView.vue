<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { ProductForm } from '@wms/ui'
import { useAuthStore } from '../stores/auth'
import { useInventoryStore } from '../stores/inventory'
import type { CreateProductRequest, UpdateProductRequest } from '@wms/api-client'

interface ProductFormData {
  sku: string; name: string; description: string; category: string; unit: string
  weight_kg: number | null; low_stock_threshold: number | null
}

const auth = useAuthStore()
const inv = useInventoryStore()

const showForm = ref(false)
const editProductId = ref<string | null>(null)
const formLoading = ref(false)
const formErrors = ref<Record<string, string>>({})

onMounted(async () => {
  if (auth.token) {
    inv.setToken(auth.token)
    await inv.fetchProducts()
  }
})

watch(() => auth.token, async (tok) => {
  if (tok) {
    inv.setToken(tok)
    await inv.fetchProducts()
  }
})

function handleCreate(): void {
  editProductId.value = null
  showForm.value = true
}

function handleEdit(productId: string): void {
  editProductId.value = productId
  showForm.value = true
}

async function handleSubmit(data: ProductFormData): Promise<void> {
  formLoading.value = true
  formErrors.value = {}

  if (editProductId.value) {
    const updateData: UpdateProductRequest = {
      name: data.name,
      description: data.description,
      category: data.category,
      unit: data.unit,
      weight_kg: data.weight_kg ?? undefined,
      low_stock_threshold: data.low_stock_threshold ?? undefined,
    }
    const result = await inv.updateProduct(editProductId.value, updateData)
    if (result) {
      showForm.value = false
      editProductId.value = null
    } else {
      formErrors.value = { _general: 'Failed to update product' }
    }
  } else {
    const createData: CreateProductRequest = {
      sku: data.sku,
      name: data.name,
      description: data.description || undefined,
      category: data.category || undefined,
      unit: data.unit || undefined,
      weight_kg: data.weight_kg ?? undefined,
      low_stock_threshold: data.low_stock_threshold ?? undefined,
    }
    const result = await inv.createProduct(createData)
    if (result) {
      showForm.value = false
    } else {
      formErrors.value = { _general: 'Failed to create product' }
    }
  }

  formLoading.value = false
}

async function handleArchive(productId: string): Promise<void> {
  if (confirm('Archive this product?')) {
    await inv.archiveProduct(productId)
  }
}

function handleCancel(): void {
  showForm.value = false
  editProductId.value = null
}

function getInitial(): Partial<ProductFormData> | undefined {
  if (!editProductId.value) return undefined
  const product = inv.products.find(p => p.id === editProductId.value)
  if (!product) return undefined
  return {
    sku: product.sku,
    name: product.name,
    description: product.description,
    category: product.category,
    unit: product.unit,
    weight_kg: product.weight_kg,
    low_stock_threshold: product.low_stock_threshold,
  }
}
</script>

<template>
  <div class="py-8 px-4 space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-semibold text-gray-800 dark:text-gray-100">Products</h1>
      <button
        v-if="!showForm"
        class="px-4 py-2 bg-blue-600 text-white rounded text-sm font-medium hover:bg-blue-700"
        @click="handleCreate"
      >
        + New Product
      </button>
    </div>

    <!-- Product Form -->
    <ProductForm
      v-if="showForm"
      :mode="editProductId ? 'edit' : 'create'"
      :initial="getInitial()"
      :loading="formLoading"
      :errors="formErrors"
      @submit="handleSubmit"
      @cancel="handleCancel"
    />

    <!-- Product List -->
    <div v-else class="bg-white border rounded-lg overflow-hidden dark:bg-gray-900 dark:border-gray-700">
      <table class="min-w-full divide-y divide-gray-200 dark:text-gray-100">
        <thead class="bg-gray-50 dark:bg-gray-800">
          <tr>
            <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase dark:text-gray-300">SKU</th>
            <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase dark:text-gray-300">Name</th>
            <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase dark:text-gray-300">Category</th>
            <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase dark:text-gray-300">Unit</th>
            <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase dark:text-gray-300">Weight</th>
            <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase dark:text-gray-300">Threshold</th>
            <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase dark:text-gray-300">Status</th>
            <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase dark:text-gray-300">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
          <tr v-if="inv.productsLoading">
            <td colspan="8" class="px-4 py-8 text-center text-gray-400">Loading...</td>
          </tr>
          <tr v-else-if="inv.products.length === 0">
            <td colspan="8" class="px-4 py-8 text-center text-gray-400">No products</td>
          </tr>
          <tr
            v-for="product in inv.products"
            :key="product.id"
            class="hover:bg-gray-50 dark:border-gray-700 dark:hover:bg-gray-800"
          >
            <td class="px-4 py-3 text-sm font-mono text-blue-600 dark:text-blue-400">{{ product.sku }}</td>
            <td class="px-4 py-3 text-sm text-gray-700 dark:text-gray-300">{{ product.name }}</td>
            <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">{{ product.category || '-' }}</td>
            <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">{{ product.unit }}</td>
            <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">{{ product.weight_kg ? product.weight_kg + ' kg' : '-' }}</td>
            <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">{{ product.low_stock_threshold }}</td>
            <td class="px-4 py-3 text-sm">
              <span
                class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium border"
                :class="product.is_active
                  ? 'bg-green-100 text-green-800 border-green-300 dark:bg-green-900/50 dark:text-green-300'
                  : 'bg-gray-100 text-gray-800 border-gray-300 dark:bg-gray-800 dark:text-gray-300'"
              >
                {{ product.is_active ? 'Active' : 'Archived' }}
              </span>
            </td>
            <td class="px-4 py-3 text-sm whitespace-nowrap">
              <button
                class="text-xs text-blue-600 hover:text-blue-800 mr-2 dark:text-blue-400 dark:hover:text-blue-300"
                @click="handleEdit(product.id)"
              >Edit</button>
              <button
                v-if="product.is_active"
                class="text-xs text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300"
                @click="handleArchive(product.id)"
              >Archive</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

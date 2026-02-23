<template>
  <n-upload
    v-model:file-list="fileList"
    list-type="image-card"
    :max="1"
    @before-upload="onBeforeUpload"
    @remove="onRemove"
  >
    <slot>
      Выберите файл
    </slot>
  </n-upload>
</template>

<script lang="ts" setup>
import { type UploadFileInfo, NUpload } from 'naive-ui'
import { onMounted, ref, watch } from 'vue'
import { useFetch } from '@/composables/useFetch'
import { type File } from '@/types'
import { useNotifications } from '@/composables/useNotifications'

const notify = useNotifications()
const fetcher = useFetch()

const value = defineModel<number | null>('value')

const { content } = defineProps<{
  content: string | null | undefined
}>()

const fileList = ref<UploadFileInfo[]>([])

const onBeforeUpload = async (data: { file: UploadFileInfo }) => {
  if (!(data.file.file instanceof File))  {
    notify.error('Файл не найден')
    return false
  }

  try {
    const response = await fetcher.uploadFile(data.file.file)

    if (response.ok) {
      fileList.value.push({
        id: response.data.uuid,
        name: response.data.origin_name,
        status: 'finished',
        url: response.data.content,
        file: data.file.file,
        type: response.data.mime_type,
      })

      value.value = response.data.id
    }
  } catch (e) {
    notify.error('При загрузке файла произошла ошибка')
    console.error(e)
  }

  return false
}

const onRemove = () => {
  value.value = null

  return false
}

const updateFileList = () => {
  if (value.value && content) {
    fileList.value.push({
      id: content,
      name: content,
      status: 'finished',
      url: content,
      // type: file.mime_type,
    })
  }
}

watch(value, (newValue) => {
  updateFileList()
  if (!newValue) {
    fileList.value = []
  }
})

// const pushFile = (file: File) => {
//   fileList.value.push({
//     id: file.uuid,
//     name: file.origin_name,
//     status: 'finished',
//     url: file.content,
//     type: file.mime_type,
//   })
//   value.value = file.id
// }
//
// defineExpose({
//   pushFile,
// })

onMounted(() => {
  updateFileList()
})
</script>

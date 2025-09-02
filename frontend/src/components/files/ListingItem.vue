<template>
  <div
    class="item"
    :class="{ 'hidden-file': name.startsWith('.') }"
    role="button"
    tabindex="0"
    :draggable="isDraggable"
    @dragstart="dragStart"
    @dragover="dragOver"
    @drop="drop"
    @click="itemClick"
    @mousedown="handleMouseDown"
    @mouseup="handleMouseUp"
    @mouseleave="handleMouseLeave"
    @touchstart="handleTouchStart"
    @touchend="handleTouchEnd"
    @touchcancel="handleTouchCancel"
    @touchmove="handleTouchMove"
    :data-dir="isDir"
    :data-type="type"
    :data-index="index"
    :aria-label="name"
    :aria-selected="isSelected"
    :data-ext="getExtension(name).toLowerCase()"
    :data-mime-type="getMimeTypeCategory(type)"
  >
    <div>
      <img
        v-if="!readOnly && type === 'image' && isThumbsEnabled"
        v-lazy="thumbnailUrl"
      />
      <i v-else class="material-icons"></i>
    </div>

    <div>
      <p class="name">{{ name }}</p>

      <p v-if="isDir" class="size dir-size" :data-order="displaySize || -1">
        <span v-if="calculatingSize" class="calculating">
          <span class="loading-dots">
            <span>计</span><span>算</span><span>中</span>
          </span>
        </span>
        <span v-else-if="displaySize !== null && displaySize >= 0" @click.stop="calculateSize" class="size-value">
          {{ filesize(displaySize) }}
        </span>
        <span v-else @click.stop="calculateSize" class="calculate-link">
          计算
        </span>
      </p>
      <p v-else class="size" :data-order="size">{{ humanSize() }}</p>

      <p class="modified">
        <time :datetime="modified">{{ humanTime() }}</time>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from "@/stores/auth";
import { useFileStore } from "@/stores/file";
import { useLayoutStore } from "@/stores/layout";

import { enableThumbs } from "@/utils/constants";
import { filesize } from "@/utils";
import dayjs from "dayjs";
import { files as api } from "@/api";
import * as upload from "@/utils/upload";
import { computed, inject, ref } from "vue";
import { useRoute, useRouter } from "vue-router";

const touches = ref<number>(0);

const longPressTimer = ref<number | null>(null);
const longPressTriggered = ref<boolean>(false);
const longPressDelay = ref<number>(500);
const startPosition = ref<{ x: number; y: number } | null>(null);
const moveThreshold = ref<number>(10);
const longPressSelected = ref<boolean>(false);

// Directory size calculation
const dirSizeValue = ref<number | null>(null);
const calculatingSize = ref<boolean>(false);

const $showError = inject<IToastError>("$showError")!;
const route = useRoute();
const router = useRouter();

const props = defineProps<{
  name: string;
  isDir: boolean;
  url: string;
  type: string;
  size: number;
  modified: string;
  index: number;
  readOnly?: boolean;
  path?: string;
}>();

const authStore = useAuthStore();
const fileStore = useFileStore();
const layoutStore = useLayoutStore();

const singleClick = computed(
  () => !props.readOnly && authStore.user?.singleClick
);
const isSelected = computed(
  () => fileStore.selected.indexOf(props.index) !== -1
);
const isDraggable = computed(
  () => !props.readOnly && authStore.user?.perm.rename
);

// Check if current view is sorted by size (which triggers auto-calculation)
const isSortedBySize = computed(() => {
  return fileStore.req?.sorting?.by === 'size';
});

// Display size: only show if explicitly calculated
const displaySize = computed(() => {
  // 1. If manually calculated locally, always show
  if (dirSizeValue.value !== null) {
    return dirSizeValue.value;
  }
  
  // 2. If sorted by size, backend calculates all sizes - show them
  if (props.isDir && isSortedBySize.value) {
    return props.size;
  }
  
  // 3. Otherwise, don't show any size (show "计算" button)
  return null;
});

const canDrop = computed(() => {
  if (!props.isDir || props.readOnly) return false;

  for (const i of fileStore.selected) {
    if (fileStore.req?.items[i].url === props.url) {
      return false;
    }
  }

  return true;
});

const thumbnailUrl = computed(() => {
  const file = {
    path: props.path,
    modified: props.modified,
  };

  return api.getPreviewURL(file as Resource, "thumb");
});

const isThumbsEnabled = computed(() => {
  return enableThumbs;
});

const humanSize = () => {
  return props.type == "invalid_link" ? "invalid link" : filesize(props.size);
};

const humanTime = () => {
  if (!props.readOnly && authStore.user?.dateFormat) {
    return dayjs(props.modified).format("L LT");
  }
  return dayjs(props.modified).fromNow();
};

const dragStart = () => {
  if (fileStore.selectedCount === 0) {
    fileStore.selected.push(props.index);
    return;
  }

  if (!isSelected.value) {
    fileStore.selected = [];
    fileStore.selected.push(props.index);
  }
};

const dragOver = (event: Event) => {
  if (!canDrop.value) return;

  event.preventDefault();
  let el = event.target as HTMLElement | null;
  if (el !== null) {
    for (let i = 0; i < 5; i++) {
      if (!el?.classList.contains("item")) {
        el = el?.parentElement ?? null;
      }
    }

    if (el !== null) el.style.opacity = "1";
  }
};

const drop = async (event: Event) => {
  if (!canDrop.value) return;
  event.preventDefault();

  if (fileStore.selectedCount === 0) return;

  let el = event.target as HTMLElement | null;
  for (let i = 0; i < 5; i++) {
    if (el !== null && !el.classList.contains("item")) {
      el = el.parentElement;
    }
  }

  const items: any[] = [];

  for (const i of fileStore.selected) {
    if (fileStore.req) {
      items.push({
        from: fileStore.req?.items[i].url,
        to: props.url + encodeURIComponent(fileStore.req?.items[i].name),
        name: fileStore.req?.items[i].name,
      });
    }
  }

  // Get url from ListingItem instance
  if (el === null) {
    return;
  }
  const path = el.__vue__.url;
  const baseItems = (await api.fetch(path)).items;

  const action = (overwrite: boolean, rename: boolean) => {
    api
      .move(items, overwrite, rename)
      .then(() => {
        fileStore.reload = true;
      })
      .catch($showError);
  };

  const conflict = upload.checkConflict(items, baseItems);

  let overwrite = false;
  let rename = false;

  if (conflict) {
    layoutStore.showHover({
      prompt: "replace-rename",
      confirm: (event: Event, option: any) => {
        overwrite = option == "overwrite";
        rename = option == "rename";

        event.preventDefault();
        layoutStore.closeHovers();
        action(overwrite, rename);
      },
    });

    return;
  }

  action(overwrite, rename);
};

const itemClick = (event: Event | KeyboardEvent) => {
  // If long press was triggered, prevent normal click behavior
  if (longPressTriggered.value) {
    longPressTriggered.value = false;
    return;
  }

  if (
    singleClick.value &&
    !(event as KeyboardEvent).ctrlKey &&
    !(event as KeyboardEvent).metaKey &&
    !(event as KeyboardEvent).shiftKey &&
    !fileStore.multiple
  )
    open();
  else click(event);
};

const click = (event: Event | KeyboardEvent) => {
  if (!singleClick.value && fileStore.selectedCount !== 0)
    event.preventDefault();

  setTimeout(() => {
    touches.value = 0;
  }, 300);

  touches.value++;
  if (touches.value > 1) {
    open();
  }

  if (fileStore.selected.indexOf(props.index) !== -1) {
    // 如果是刚刚长按选中的文件，在保护期内不取消选中
    if (longPressSelected.value) {
      return;
    }
    fileStore.removeSelected(props.index);
    return;
  }

  if ((event as KeyboardEvent).shiftKey && fileStore.selected.length > 0) {
    let fi = 0;
    let la = 0;

    if (props.index > fileStore.selected[0]) {
      fi = fileStore.selected[0] + 1;
      la = props.index;
    } else {
      fi = props.index;
      la = fileStore.selected[0] - 1;
    }

    for (; fi <= la; fi++) {
      if (fileStore.selected.indexOf(fi) == -1) {
        fileStore.selected.push(fi);
      }
    }

    return;
  }

  if (
    !singleClick.value &&
    !(event as KeyboardEvent).ctrlKey &&
    !(event as KeyboardEvent).metaKey &&
    !fileStore.multiple
  ) {
    fileStore.selected = [];
  }
  fileStore.selected.push(props.index);
};

const open = () => {
  // 在打开文件/目录前保存当前滚动位置
  const currentPath = fileStore.req?.path || route.path;
  fileStore.saveScrollPosition(currentPath);

  router.push({ path: props.url });
};

const getExtension = (fileName: string): string => {
  const lastDotIndex = fileName.lastIndexOf(".");
  if (lastDotIndex === -1) {
    return fileName;
  }
  return fileName.substring(lastDotIndex);
};

const getMimeTypeCategory = (fileType: string) => {
  // Return a category based on the type to override CSS icon rules
  // This will be used for data-mime-type attribute
  // Only return known types that should override default extension-based icons
  if (fileType === "text") return "text";
  if (fileType === "image") return "image";
  if (fileType === "audio") return "audio";
  if (fileType === "video") return "video";
  if (fileType === "pdf") return "pdf";
  if (fileType === "archive") return "archive";
  if (fileType === "code") return "code";
  // Return empty string for unknown types so default extension rules apply
  return "";
};

const calculateSize = async () => {
  if (!props.isDir || calculatingSize.value) return;
  
  calculatingSize.value = true;
  try {
    // Use fetchURL directly for the API call
    const { fetchURL, removePrefix } = await import("@/api/utils");
    const cleanPath = removePrefix(props.url);
    const response = await fetchURL(`/api/dirsize${cleanPath}`, {});
    const data = await response.json();
    dirSizeValue.value = data.size;
    
    // Update the item in the store if it exists
    if (fileStore.req && fileStore.req.items && fileStore.req.items[props.index]) {
      fileStore.req.items[props.index].size = data.size;
    }
  } catch (e: any) {
    $showError(e);
  } finally {
    calculatingSize.value = false;
  }
};

// Long-press helper functions
const startLongPress = (clientX: number, clientY: number) => {
  startPosition.value = { x: clientX, y: clientY };
  longPressTimer.value = window.setTimeout(() => {
    handleLongPress();
  }, longPressDelay.value);
};

const cancelLongPress = () => {
  if (longPressTimer.value !== null) {
    window.clearTimeout(longPressTimer.value);
    longPressTimer.value = null;
  }
  startPosition.value = null;
};

const handleLongPress = () => {
  if (singleClick.value) {
    longPressTriggered.value = true;
    // 长按选中逻辑：如果未选中则添加到选中列表，如果已选中则保持选中状态
    if (!isSelected.value) {
      if (!fileStore.multiple) {
        fileStore.selected = [];
      }
      fileStore.selected.push(props.index);
      // 标记这个文件是通过长按选中的
      longPressSelected.value = true;
      // 设置保护期，避免立即被取消选中
      setTimeout(() => {
        longPressSelected.value = false;
      }, 500);
    }
    // 长按后自动开启多选模式，保持逻辑一致性
    if (!fileStore.multiple) {
      fileStore.multiple = true;
    }
    // 长按时不执行普通的click逻辑，避免取消选中
  }
  cancelLongPress();
};

const checkMovement = (clientX: number, clientY: number): boolean => {
  if (!startPosition.value) return false;

  const deltaX = Math.abs(clientX - startPosition.value.x);
  const deltaY = Math.abs(clientY - startPosition.value.y);

  return deltaX > moveThreshold.value || deltaY > moveThreshold.value;
};

// Event handlers
const handleMouseDown = (event: MouseEvent) => {
  if (event.button === 0) {
    startLongPress(event.clientX, event.clientY);
  }
};

const handleMouseUp = () => {
  cancelLongPress();
  // 延迟重置长按标志，允许后续的点击事件正常处理
  setTimeout(() => {
    longPressTriggered.value = false;
  }, 50);
};

const handleMouseLeave = () => {
  cancelLongPress();
  // 重置长按标志
  longPressTriggered.value = false;
};

const handleTouchStart = (event: TouchEvent) => {
  if (event.touches.length === 1) {
    const touch = event.touches[0];
    startLongPress(touch.clientX, touch.clientY);
  }
};

const handleTouchEnd = () => {
  cancelLongPress();
  // 延迟重置长按标志，允许后续的点击事件正常处理
  setTimeout(() => {
    longPressTriggered.value = false;
  }, 50);
};

const handleTouchCancel = () => {
  cancelLongPress();
  // 重置长按标志
  longPressTriggered.value = false;
};

const handleTouchMove = (event: TouchEvent) => {
  if (event.touches.length === 1 && startPosition.value) {
    const touch = event.touches[0];
    if (checkMovement(touch.clientX, touch.clientY)) {
      cancelLongPress();
    }
  }
};
</script>

<style scoped>
.dir-size {
  position: relative;
}

/* Default state - use blue color */
.calculate-link {
  cursor: pointer;
  color: var(--blue);
  transition: all 0.2s;
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 0.85em;
}

.calculate-link:hover {
  background: rgba(192, 132, 252, 0.1);
  color: var(--dark-blue);
}

/* When item is selected, inherit the text color (white on blue background) */
.item[aria-selected="true"] .dir-size .calculate-link {
  color: var(--iconSecondary) !important;
  opacity: 0.9;
}

.item[aria-selected="true"] .dir-size .calculate-link:hover {
  background: rgba(255, 255, 255, 0.15);
  opacity: 1;
  color: var(--iconSecondary) !important;
}

/* Size value styles */
.size-value {
  cursor: pointer;
  transition: color 0.2s;
}

.size-value:hover {
  color: var(--blue);
}

.item[aria-selected="true"] .dir-size .size-value {
  color: var(--iconSecondary) !important;
}

.item[aria-selected="true"] .dir-size .size-value:hover {
  opacity: 0.8;
  color: var(--iconSecondary) !important;
}

/* Calculating state */
.calculating {
  display: inline-flex;
  align-items: center;
  color: var(--blue);
}

.item[aria-selected="true"] .dir-size .calculating {
  color: var(--iconSecondary) !important;
}

.item[aria-selected="true"] .dir-size .loading-dots {
  color: var(--iconSecondary) !important;
}

.loading-dots {
  display: inline-flex;
  font-size: 0.85em;
}

.loading-dots span {
  animation: wave 1.4s ease-in-out infinite;
  display: inline-block;
}

.loading-dots span:nth-child(1) {
  animation-delay: 0s;
}

.loading-dots span:nth-child(2) {
  animation-delay: 0.1s;
}

.loading-dots span:nth-child(3) {
  animation-delay: 0.2s;
}

@keyframes wave {
  0%, 60%, 100% {
    transform: translateY(0);
    opacity: 1;
  }
  30% {
    transform: translateY(-3px);
    opacity: 0.7;
  }
}

/* Hidden files (starting with .) appear with reduced opacity */
.hidden-file {
  opacity: 0.6;
}

.hidden-file:hover {
  opacity: 0.8;
}

/* When selected, hidden files show full opacity */
.hidden-file.item[aria-selected="true"] {
  opacity: 1;
}
</style>

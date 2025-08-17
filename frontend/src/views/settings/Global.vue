<template>
  <errors v-if="error" :errorCode="error.status" />
  <div class="row" v-else-if="!layoutStore.loading && settings !== null">
    <div class="column">
      <form class="card" @submit.prevent="save">
        <div class="card-title">
          <h2>{{ t("settings.globalSettings") }}</h2>
        </div>

        <div class="card-content">
          <p>
            <input type="checkbox" v-model="settings.signup" />
            {{ t("settings.allowSignup") }}
          </p>

          <p>
            <input type="checkbox" v-model="settings.createUserDir" />
            {{ t("settings.createUserDir") }}
          </p>

          <p>
            <label class="small">{{ t("settings.userHomeBasePath") }}</label>
            <input
              class="input input--block"
              type="text"
              v-model="settings.userHomeBasePath"
            />
          </p>

          <p>
            <label for="minimumPasswordLength">{{
              t("settings.minimumPasswordLength")
            }}</label>
            <vue-number-input
              controls
              v-model.number="settings.minimumPasswordLength"
              id="minimumPasswordLength"
              :min="1"
            />
          </p>

          <h3>{{ t("settings.rules") }}</h3>
          <p class="small">{{ t("settings.globalRules") }}</p>
          <rules v-model:rules="settings.rules" />

          <div v-if="enableExec">
            <h3>{{ t("settings.executeOnShell") }}</h3>
            <p class="small">{{ t("settings.executeOnShellDescription") }}</p>
            <input
              class="input input--block"
              type="text"
              placeholder="bash -c, cmd /c, ..."
              v-model="shellValue"
            />
          </div>

          <h3>{{ t("settings.branding") }}</h3>

          <i18n-t
            keypath="settings.brandingHelp"
            tag="p"
            class="small"
            scope="global"
          >
            <a
              class="link"
              target="_blank"
              href="https://filebrowser.org/configuration#custom-branding"
              >{{ t("settings.documentation") }}</a
            >
          </i18n-t>

          <p>
            <input
              type="checkbox"
              v-model="settings.branding.disableExternal"
              id="branding-links"
            />
            {{ t("settings.disableExternalLinks") }}
          </p>

          <p>
            <input
              type="checkbox"
              v-model="settings.branding.disableUsedPercentage"
              id="branding-used-disk"
            />
            {{ t("settings.disableUsedDiskPercentage") }}
          </p>

          <p>
            <label for="theme">{{ t("settings.themes.title") }}</label>
            <themes
              class="input input--block"
              v-model:theme="settings.branding.theme"
              id="theme"
            ></themes>
          </p>

          <p>
            <label for="branding-name">{{ t("settings.instanceName") }}</label>
            <input
              class="input input--block"
              type="text"
              v-model="settings.branding.name"
              id="branding-name"
            />
          </p>

          <p>
            <label for="branding-files">{{
              t("settings.brandingDirectoryPath")
            }}</label>
            <input
              class="input input--block"
              type="text"
              v-model="settings.branding.files"
              id="branding-files"
            />
          </p>

          <h3>{{ t("settings.tusUploads") }}</h3>

          <p class="small">{{ t("settings.tusUploadsHelp") }}</p>

          <div class="tusConditionalSettings">
            <p>
              <label for="tus-chunkSize">{{
                t("settings.tusUploadsChunkSize")
              }}</label>
              <input
                class="input input--block"
                type="text"
                v-model="formattedChunkSize"
                id="tus-chunkSize"
              />
            </p>

            <p>
              <label for="tus-retryCount">{{
                t("settings.tusUploadsRetryCount")
              }}</label>
              <vue-number-input
                controls
                v-model.number="settings.tus.retryCount"
                id="tus-retryCount"
                :min="0"
              />
            </p>
          </div>
        </div>

        <div class="card-action">
          <input
            class="button button--flat"
            type="submit"
            :value="t('buttons.update')"
          />
        </div>
      </form>
    </div>

    <div class="column">
      <form class="card" @submit.prevent="save">
        <div class="card-title">
          <h2>{{ t("settings.userDefaults") }}</h2>
        </div>

        <div class="card-content">
          <p class="small">{{ t("settings.defaultUserDescription") }}</p>

          <user-form
            :isNew="false"
            :isDefault="true"
            v-model:user="settings.defaults"
          />
        </div>

        <div class="card-action">
          <input
            class="button button--flat"
            type="submit"
            :value="t('buttons.update')"
          />
        </div>
      </form>
    </div>

    <div class="column">
      <form v-if="enableExec" class="card" @submit.prevent="save">
        <div class="card-title">
          <h2>{{ t("settings.commandRunner") }}</h2>
        </div>

        <div class="card-content">
          <i18n-t
            keypath="settings.commandRunnerHelp"
            tag="p"
            class="small"
            scope="global"
          >
            <code>FILE</code>
            <code>SCOPE</code>
            <a
              class="link"
              target="_blank"
              href="https://filebrowser.org/configuration.html#command-runner"
              >{{ t("settings.documentation") }}</a
            >
          </i18n-t>

          <div
            v-for="(command, key) in settings.commands"
            :key="key"
            class="collapsible"
          >
            <input :id="key" type="checkbox" />
            <label :for="key">
              <p>{{ capitalize(key) }}</p>
              <i class="material-icons">arrow_drop_down</i>
            </label>
            <div class="collapse">
              <textarea
                class="input input--block input--textarea"
                v-model.trim="commandObject[key]"
              ></textarea>
            </div>
          </div>
        </div>

        <div class="card-action">
          <input
            class="button button--flat"
            type="submit"
            :value="t('buttons.update')"
          />
        </div>
      </form>
    </div>

    <div class="column">
      <form class="card" @submit.prevent="save">
        <div class="card-title">
          <h2>{{ t("settings.mimeTypes") }}</h2>
        </div>

        <div class="card-content">
          <p class="small">{{ t("settings.mimeTypesHelp") }}</p>

          <div class="mime-types-section">
            <h4>{{ t("settings.extensionMimeTypes") }}</h4>
            <div class="mime-entries-container" v-if="Object.keys(mimeTypes).length > 0">
              <div
                v-for="(mime, ext) in mimeTypes"
                :key="ext"
                class="mime-entry"
              >
                <div class="mime-entry-inputs">
                  <input
                    class="input mime-ext-input"
                    type="text"
                    :value="ext"
                    @change="updateMimeExtension($event, String(ext), mime)"
                    @keydown.enter.prevent
                    placeholder=".ext"
                    spellcheck="false"
                  />
                  <div class="mime-type-wrapper">
                    <select
                      :value="mime"
                      @change="updateMimeTypeValue(String(ext), $event)"
                      class="input mime-type-select"
                    >
                      <option value="">
                        {{ t("settings.customMimeType") }}
                      </option>
                      <option
                        v-for="mimeType in commonMimeTypes"
                        :key="mimeType.value"
                        :value="mimeType.value"
                      >
                        {{ mimeType.label }} ({{ mimeType.value }})
                      </option>
                    </select>
                    <input
                      v-if="!commonMimeTypes.find((m) => m.value === mime)"
                      class="input mime-type-input"
                      type="text"
                      :value="mime"
                      @input="updateMimeTypeValue(String(ext), $event)"
                      placeholder="application/octet-stream"
                      spellcheck="false"
                    />
                  </div>
                </div>
                <button
                  type="button"
                  class="button button--icon button--flat mime-delete-btn"
                  @click="removeMimeType(String(ext))"
                  :title="t('buttons.delete')"
                >
                  <i class="material-icons">clear</i>
                </button>
              </div>
            </div>
            <div v-else class="mime-empty-state">
              <i class="material-icons">folder_open</i>
              <p>{{ t('settings.noCustomMimeTypes', '没有自定义的 MIME 类型') }}</p>
            </div>
            <button
              type="button"
              class="button button--small button--outlined mime-add-btn"
              @click="addMimeType"
            >
              <i class="material-icons">add</i>
              <span>{{ t("settings.addMimeType") }}</span>
            </button>
          </div>

          <div class="mime-types-section">
            <h4>{{ t("settings.filenameMimeTypes") }}</h4>
            <div class="mime-entries-container" v-if="Object.keys(filenameMimes).length > 0">
              <div
                v-for="(mime, filename) in filenameMimes"
                :key="filename"
                class="mime-entry"
              >
                <div class="mime-entry-inputs">
                  <input
                    class="input mime-ext-input"
                    type="text"
                    :value="filename"
                    @change="updateFilenameMime($event, String(filename), mime)"
                    @keydown.enter.prevent
                    placeholder="README.md"
                    spellcheck="false"
                  />
                  <div class="mime-type-wrapper">
                    <select
                      :value="mime"
                      @change="
                        updateFilenameMimeValue(String(filename), $event)
                      "
                      class="input mime-type-select"
                    >
                      <option value="">
                        {{ t("settings.customMimeType") }}
                      </option>
                      <option
                        v-for="mimeType in commonMimeTypes"
                        :key="mimeType.value"
                        :value="mimeType.value"
                      >
                        {{ mimeType.label }} ({{ mimeType.value }})
                      </option>
                    </select>
                    <input
                      v-if="!commonMimeTypes.find((m) => m.value === mime)"
                      class="input mime-type-input"
                      type="text"
                      :value="mime"
                      @input="updateFilenameMimeValue(String(filename), $event)"
                      placeholder="application/octet-stream"
                      spellcheck="false"
                    />
                  </div>
                </div>
                <button
                  type="button"
                  class="button button--icon button--flat mime-delete-btn"
                  @click="removeFilenameMime(String(filename))"
                  :title="t('buttons.delete')"
                >
                  <i class="material-icons">clear</i>
                </button>
              </div>
            </div>
            <div v-else class="mime-empty-state">
              <i class="material-icons">insert_drive_file</i>
              <p>{{ t('settings.noFilenameMimeTypes', '没有特定文件名的 MIME 类型') }}</p>
            </div>
            <button
              type="button"
              class="button button--small button--outlined mime-add-btn"
              @click="addFilenameMime"
            >
              <i class="material-icons">add</i>
              <span>{{ t("settings.addFilenameMime") }}</span>
            </button>
          </div>
        </div>

        <div class="card-action">
          <input
            class="button button--flat"
            type="submit"
            :value="t('buttons.update')"
          />
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { settings as api } from "@/api";
import { StatusError } from "@/api/utils";
import Rules from "@/components/settings/Rules.vue";
import Themes from "@/components/settings/Themes.vue";
import UserForm from "@/components/settings/UserForm.vue";
import { useLayoutStore } from "@/stores/layout";
import { enableExec } from "@/utils/constants";
import { getTheme, setTheme } from "@/utils/theme";
import Errors from "@/views/Errors.vue";
import { computed, inject, onBeforeUnmount, onMounted, ref } from "vue";
import { useI18n } from "vue-i18n";

const error = ref<StatusError | null>(null);
const originalSettings = ref<ISettings | null>(null);
const settings = ref<ISettings | null>(null);
const debounceTimeout = ref<number | null>(null);

const commandObject = ref<{
  [key: string]: string[] | string;
}>({});
const shellValue = ref<string>("");

const $showError = inject<IToastError>("$showError")!;
const $showSuccess = inject<IToastSuccess>("$showSuccess")!;

const { t } = useI18n();

const layoutStore = useLayoutStore();

// Computed properties for safe access
const mimeTypes = computed(() => settings.value?.mimeTypes || {});
const filenameMimes = computed(() => settings.value?.filenameMimes || {});

// Common MIME types with Chinese descriptions
const commonMimeTypes = [
  { label: "纯文本", value: "text/plain" },
  { label: "HTML网页", value: "text/html" },
  { label: "CSS样式表", value: "text/css" },
  { label: "JavaScript脚本", value: "text/javascript" },
  { label: "JSON数据", value: "application/json" },
  { label: "XML文档", value: "text/xml" },
  { label: "Markdown文档", value: "text/markdown" },
  { label: "YAML配置", value: "text/yaml" },
  { label: "CSV表格", value: "text/csv" },
  { label: "PDF文档", value: "application/pdf" },
  { label: "Word文档", value: "application/msword" },
  { label: "Excel表格", value: "application/vnd.ms-excel" },
  { label: "PowerPoint演示", value: "application/vnd.ms-powerpoint" },
  { label: "ZIP压缩包", value: "application/zip" },
  { label: "RAR压缩包", value: "application/x-rar-compressed" },
  { label: "7Z压缩包", value: "application/x-7z-compressed" },
  { label: "TAR归档", value: "application/x-tar" },
  { label: "GZIP压缩", value: "application/gzip" },
  { label: "PNG图片", value: "image/png" },
  { label: "JPEG图片", value: "image/jpeg" },
  { label: "GIF图片", value: "image/gif" },
  { label: "SVG矢量图", value: "image/svg+xml" },
  { label: "MP3音频", value: "audio/mpeg" },
  { label: "WAV音频", value: "audio/wav" },
  { label: "OGG音频", value: "audio/ogg" },
  { label: "MP4视频", value: "video/mp4" },
  { label: "WebM视频", value: "video/webm" },
  { label: "AVI视频", value: "video/x-msvideo" },
  { label: "二进制文件", value: "application/octet-stream" },
];

const formattedChunkSize = computed({
  get() {
    return settings?.value?.tus?.chunkSize
      ? formatBytes(settings?.value?.tus?.chunkSize)
      : "";
  },
  set(value: string) {
    // Use debouncing to allow the user to type freely without
    // interruption by the formatter
    // Clear the previous timeout if it exists
    if (debounceTimeout.value) {
      clearTimeout(debounceTimeout.value);
    }

    // Set a new timeout to apply the format after a short delay
    debounceTimeout.value = window.setTimeout(() => {
      if (settings.value) settings.value.tus.chunkSize = parseBytes(value);
    }, 1500);
  },
});

// Define funcs
const capitalize = (name: string, where: string | RegExp = "_") => {
  if (where === "caps") where = /(?=[A-Z])/;
  const split = name.split(where);
  name = "";

  for (let i = 0; i < split.length; i++) {
    name += split[i].charAt(0).toUpperCase() + split[i].slice(1) + " ";
  }

  return name.slice(0, -1);
};

const save = async () => {
  if (settings.value === null) return false;
  const newSettings: ISettings = {
    ...settings.value,
    shell:
      settings.value?.shell
        .join(" ")
        .trim()
        .split(" ")
        .filter((s: string) => s !== "") ?? [],
    commands: {},
  };

  const keys = Object.keys(settings.value.commands) as Array<
    keyof SettingsCommand
  >;
  for (const key of keys) {
    // not sure if we can safely assume non-null
    const newValue = commandObject.value[key];
    if (!newValue) continue;

    if (Array.isArray(newValue)) {
      newSettings.commands[key] = newValue;
    } else if (key in commandObject.value) {
      newSettings.commands[key] = newValue
        .split("\n")
        .filter((cmd: string) => cmd !== "");
    }
  }
  newSettings.shell = shellValue.value
    .trim()
    .split(" ")
    .filter((s) => s !== "");

  if (newSettings.branding.theme !== getTheme()) {
    setTheme(newSettings.branding.theme);
  }

  try {
    await api.update(newSettings);
    $showSuccess(t("settings.settingsUpdated"));
  } catch (e: any) {
    $showError(e);
  }

  return true;
};
// Parse the user-friendly input (e.g., "20M" or "1T") to bytes
const parseBytes = (input: string) => {
  const regex = /^(\d+)(\.\d+)?(B|K|KB|M|MB|G|GB|T|TB)?$/i;
  const matches = input.match(regex);
  if (matches) {
    const size = parseFloat(matches[1].concat(matches[2] || ""));
    let unit: keyof SettingsUnit =
      matches[3].toUpperCase() as keyof SettingsUnit;
    if (!unit.endsWith("B")) {
      unit += "B";
    }
    const units: SettingsUnit = {
      KB: 1024,
      MB: 1024 ** 2,
      GB: 1024 ** 3,
      TB: 1024 ** 4,
    };
    return size * (units[unit as keyof SettingsUnit] || 1);
  } else {
    return 1024 ** 2;
  }
};
// Format the chunk size in bytes to user-friendly format
const formatBytes = (bytes: number) => {
  const units = ["B", "KB", "MB", "GB", "TB"];
  let size = bytes;
  let unitIndex = 0;
  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024;
    unitIndex++;
  }
  return `${size}${units[unitIndex]}`;
};

// MIME Type management methods
const addMimeType = () => {
  if (!settings.value) return;
  if (!settings.value.mimeTypes) {
    settings.value.mimeTypes = {};
  }
  // Add empty entry that user can fill
  let counter = 1;
  let ext = `.ext${counter}`;
  while (settings.value.mimeTypes[ext]) {
    counter++;
    ext = `.ext${counter}`;
  }
  settings.value.mimeTypes[ext] = "text/plain";
};

const removeMimeType = (ext: string) => {
  if (!settings.value?.mimeTypes) return;
  delete settings.value.mimeTypes[ext];
};

const updateMimeExtension = (event: Event, oldExt: string, mime: string) => {
  const target = event.target as HTMLInputElement;
  const newExt = target.value.trim();
  if (!settings.value?.mimeTypes || newExt === oldExt) return;

  // Ensure extension starts with a dot
  const normalizedExt = newExt.startsWith(".") ? newExt : `.${newExt}`;

  delete settings.value.mimeTypes[oldExt];
  settings.value.mimeTypes[normalizedExt] = mime;
};

const addFilenameMime = () => {
  if (!settings.value) return;
  if (!settings.value.filenameMimes) {
    settings.value.filenameMimes = {};
  }
  // Add empty entry that user can fill
  let counter = 1;
  let filename = `file${counter}.txt`;
  while (settings.value.filenameMimes[filename]) {
    counter++;
    filename = `file${counter}.txt`;
  }
  settings.value.filenameMimes[filename] = "text/plain";
};

const removeFilenameMime = (filename: string) => {
  if (!settings.value?.filenameMimes) return;
  delete settings.value.filenameMimes[filename];
};

const updateFilenameMime = (
  event: Event,
  oldFilename: string,
  mime: string
) => {
  const target = event.target as HTMLInputElement;
  const newFilename = target.value.trim();
  if (!settings.value?.filenameMimes || newFilename === oldFilename) return;

  delete settings.value.filenameMimes[oldFilename];
  settings.value.filenameMimes[newFilename] = mime;
};

// Update MIME type value
const updateMimeTypeValue = (ext: string, event: Event) => {
  if (!settings.value?.mimeTypes) return;
  const target = event.target as HTMLInputElement | HTMLSelectElement;
  settings.value.mimeTypes[ext] = target.value;
};

// Update filename MIME value
const updateFilenameMimeValue = (filename: string, event: Event) => {
  if (!settings.value?.filenameMimes) return;
  const target = event.target as HTMLInputElement | HTMLSelectElement;
  settings.value.filenameMimes[filename] = target.value;
};

// Define Hooks

onMounted(async () => {
  try {
    layoutStore.loading = true;
    const original: ISettings = await api.get();
    const newSettings: ISettings = { ...original, commands: {} };

    const keys = Object.keys(original.commands) as Array<keyof SettingsCommand>;
    for (const key of keys) {
      newSettings.commands[key] = original.commands[key];
      commandObject.value[key] = original.commands[key]!.join("\n");
    }

    // Initialize MIME types if they don't exist
    if (!newSettings.mimeTypes) {
      newSettings.mimeTypes = {};
    }
    if (!newSettings.filenameMimes) {
      newSettings.filenameMimes = {};
    }

    originalSettings.value = original;
    settings.value = newSettings;
    shellValue.value = newSettings.shell.join(" ");
  } catch (err) {
    if (err instanceof Error) {
      error.value = err;
    }
  } finally {
    layoutStore.loading = false;
  }
});

// Clear the debounce timeout when the component is destroyed
onBeforeUnmount(() => {
  if (debounceTimeout.value) {
    clearTimeout(debounceTimeout.value);
  }
});
</script>

<style scoped>
.mime-types-section {
  margin-bottom: 2rem;
}

.mime-types-section h4 {
  margin-bottom: 1rem;
  font-size: 0.95rem;
  font-weight: 500;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.mime-entries-container {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-bottom: 1rem;
  padding: 0.75rem;
  background: var(--background-secondary);
  border-radius: 8px;
  border: 1px solid var(--border-color);
}

.mime-empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem 1rem;
  margin-bottom: 1rem;
  background: var(--background-secondary);
  border-radius: 8px;
  border: 1px dashed var(--border-color);
  color: var(--text-tertiary);
}

.mime-empty-state i {
  font-size: 3rem;
  margin-bottom: 0.5rem;
  opacity: 0.5;
}

.mime-empty-state p {
  margin: 0;
  font-size: 0.9rem;
}

.mime-entry {
  display: flex;
  gap: 0.75rem;
  align-items: center;
  background: var(--background);
  padding: 0.75rem;
  border-radius: 6px;
  transition: all 0.2s ease;
  border: 1px solid transparent;
}

.mime-entry:hover {
  border-color: var(--primary);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.mime-entry-inputs {
  display: flex;
  gap: 0.75rem;
  flex: 1;
  align-items: center;
  flex-wrap: wrap;
}

.mime-ext-input {
  flex: 0 0 180px;
  padding: 0.5rem 0.75rem;
  font-size: 0.9rem;
  font-family: "Monaco", "Menlo", "Courier New", monospace;
  background: var(--background-secondary);
  border: 1px solid var(--border-color);
  border-radius: 4px;
  transition: all 0.2s ease;
}

.mime-ext-input:focus {
  border-color: var(--primary);
  background: var(--background);
  outline: none;
  box-shadow: 0 0 0 3px rgba(var(--primary-rgb), 0.1);
}

.mime-type-wrapper {
  flex: 1;
  min-width: 300px;
  position: relative;
}

.mime-type-select {
  width: 100%;
  padding: 0.5rem 0.75rem;
  padding-right: 2rem;
  font-size: 0.9rem;
  background: var(--background-secondary);
  border: 1px solid var(--border-color);
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s ease;
  appearance: none;
  background-image: url("data:image/svg+xml;charset=UTF-8,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3e%3cpolyline points='6 9 12 15 18 9'%3e%3c/polyline%3e%3c/svg%3e");
  background-repeat: no-repeat;
  background-position: right 0.5rem center;
  background-size: 1.2rem;
}

.mime-type-select:focus {
  border-color: var(--primary);
  background-color: var(--background);
  outline: none;
  box-shadow: 0 0 0 3px rgba(var(--primary-rgb), 0.1);
}

.mime-type-select option {
  padding: 0.5rem;
  background: var(--background);
  color: var(--text);
}

.mime-type-input {
  width: 100%;
  padding: 0.5rem 0.75rem;
  font-size: 0.9rem;
  font-family: "Monaco", "Menlo", "Courier New", monospace;
  background: var(--background-secondary);
  border: 1px solid var(--border-color);
  border-radius: 4px;
  transition: all 0.2s ease;
}

.mime-type-input:focus {
  border-color: var(--primary);
  background: var(--background);
  outline: none;
  box-shadow: 0 0 0 3px rgba(var(--primary-rgb), 0.1);
}

.mime-delete-btn {
  flex: 0 0 auto;
  width: 32px;
  height: 32px;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: all 0.2s ease;
  color: var(--text-tertiary);
}

.mime-delete-btn:hover {
  background: var(--danger-light);
  color: var(--danger);
  transform: scale(1.1);
}

.mime-delete-btn i {
  font-size: 18px;
}

.mime-add-btn {
  align-self: flex-start;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  border: 1px dashed var(--border-color);
  background: transparent;
  color: var(--primary);
  transition: all 0.2s ease;
}

.mime-add-btn:hover {
  border-color: var(--primary);
  background: rgba(var(--primary-rgb), 0.05);
  transform: translateY(-1px);
}

.mime-add-btn i {
  font-size: 18px;
}

/* Dark mode specific adjustments */
.dark .mime-entry {
  background: var(--background-dark);
}

.dark .mime-entries-container {
  background: var(--background-dark-secondary);
}

.dark .mime-type-select,
.dark .mime-type-input,
.dark .mime-ext-input {
  background: var(--background-dark);
  border-color: var(--border-dark);
}

.dark .mime-type-select:focus,
.dark .mime-type-input:focus,
.dark .mime-ext-input:focus {
  background: var(--background-dark-elevated);
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .mime-entry-inputs {
    flex-direction: column;
    align-items: stretch;
  }

  .mime-ext-input {
    flex: 1;
  }

  .mime-type-wrapper {
    min-width: auto;
  }
}
</style>

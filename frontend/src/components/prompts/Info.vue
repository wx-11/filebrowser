<template>
  <div class="card floating">
    <div class="card-title">
      <h2>{{ $t("prompts.fileInfo") }}</h2>
    </div>

    <div class="card-content">
      <p v-if="selected.length > 1">
        {{ $t("prompts.filesSelected", { count: selected.length }) }}
      </p>

      <p class="break-word" v-if="selected.length < 2">
        <strong>{{ $t("prompts.displayName") }}</strong> {{ name }}
      </p>

      <p v-if="!dir || selected.length > 1">
        <strong>{{ $t("prompts.size") }}:</strong>
        <span id="content_length"></span> {{ humanSize }}
        <span v-if="!dir && sizeInBytes" class="size-bytes"> ({{ sizeInBytes.toLocaleString() }} bytes)</span>
      </p>
      
      <p v-if="dir && selected.length === 1">
        <strong>{{ $t("prompts.size") }}:</strong>
        <span v-if="calculatingDirSize" class="calculating-size">
          <i class="material-icons spin">autorenew</i> {{ $t("prompts.calculating") }}...
        </span>
        <span v-else-if="dirSize !== null">
          {{ humanDirSize }}
          <span class="size-bytes"> ({{ dirSize.toLocaleString() }} bytes)</span>
        </span>
        <button v-else @click="calculateDirSize" class="calculate-size-btn">
          {{ $t("prompts.calculateSize") }}
        </button>
      </p>

      <div v-if="resolution">
        <strong>{{ $t("prompts.resolution") }}:</strong>
        {{ resolution.width }} x {{ resolution.height }}
      </div>

      <p v-if="selected.length < 2" :title="modTime">
        <strong>{{ $t("prompts.lastModified") }}:</strong> {{ humanTime }}
      </p>

      <template v-if="dir && selected.length === 0">
        <p>
          <strong>{{ $t("prompts.numberFiles") }}:</strong> {{ req.numFiles }}
        </p>
        <p>
          <strong>{{ $t("prompts.numberDirs") }}:</strong> {{ req.numDirs }}
        </p>
      </template>

      <template v-if="!dir">
        <p>
          <strong>MD5: </strong
          ><code
            ><a
              @click="checksum($event, 'md5')"
              @keypress.enter="checksum($event, 'md5')"
              tabindex="2"
              >{{ $t("prompts.show") }}</a
            ></code
          >
        </p>
        <p>
          <strong>SHA1: </strong
          ><code
            ><a
              @click="checksum($event, 'sha1')"
              @keypress.enter="checksum($event, 'sha1')"
              tabindex="3"
              >{{ $t("prompts.show") }}</a
            ></code
          >
        </p>
        <p>
          <strong>SHA256: </strong
          ><code
            ><a
              @click="checksum($event, 'sha256')"
              @keypress.enter="checksum($event, 'sha256')"
              tabindex="4"
              >{{ $t("prompts.show") }}</a
            ></code
          >
        </p>
        <p>
          <strong>SHA512: </strong
          ><code
            ><a
              @click="checksum($event, 'sha512')"
              @keypress.enter="checksum($event, 'sha512')"
              tabindex="5"
              >{{ $t("prompts.show") }}</a
            ></code
          >
        </p>
      </template>
    </div>

    <div class="card-action">
      <button
        id="focus-prompt"
        type="submit"
        @click="closeHovers"
        class="button button--flat"
        :aria-label="$t('buttons.ok')"
        :title="$t('buttons.ok')"
      >
        {{ $t("buttons.ok") }}
      </button>
    </div>
  </div>
</template>

<script>
import { mapActions, mapState } from "pinia";
import { useFileStore } from "@/stores/file";
import { useLayoutStore } from "@/stores/layout";
import { filesize } from "@/utils";
import dayjs from "dayjs";
import { files as api } from "@/api";

export default {
  name: "info",
  inject: ["$showError"],
  data() {
    return {
      dirSize: null,
      calculatingDirSize: false,
    };
  },
  computed: {
    ...mapState(useFileStore, [
      "req",
      "selected",
      "selectedCount",
      "isListing",
    ]),
    humanSize: function () {
      if (this.selectedCount === 0 || !this.isListing) {
        return filesize(this.req.size);
      }

      let sum = 0;

      for (const selected of this.selected) {
        sum += this.req.items[selected].size;
      }

      return filesize(sum);
    },
    sizeInBytes: function () {
      if (this.selectedCount === 0 || !this.isListing) {
        return this.req.size;
      }

      let sum = 0;
      for (const selected of this.selected) {
        sum += this.req.items[selected].size;
      }
      return sum;
    },
    humanDirSize: function () {
      return this.dirSize !== null ? filesize(this.dirSize) : "";
    },
    humanTime: function () {
      if (this.selectedCount === 0) {
        return dayjs(this.req.modified).fromNow();
      }

      return dayjs(this.req.items[this.selected[0]].modified).fromNow();
    },
    modTime: function () {
      if (this.selectedCount === 0) {
        return new Date(Date.parse(this.req.modified)).toLocaleString();
      }

      return new Date(
        Date.parse(this.req.items[this.selected[0]].modified)
      ).toLocaleString();
    },
    name: function () {
      return this.selectedCount === 0
        ? this.req.name
        : this.req.items[this.selected[0]].name;
    },
    dir: function () {
      return (
        this.selectedCount > 1 ||
        (this.selectedCount === 0
          ? this.req.isDir
          : this.req.items[this.selected[0]].isDir)
      );
    },
    resolution: function () {
      if (this.selectedCount === 1) {
        const selectedItem = this.req.items[this.selected[0]];
        if (selectedItem && selectedItem.type === "image") {
          return selectedItem.resolution;
        }
      } else if (this.req && this.req.type === "image") {
        return this.req.resolution;
      }
      return null;
    },
  },
  methods: {
    ...mapActions(useLayoutStore, ["closeHovers"]),
    checksum: async function (event, algo) {
      event.preventDefault();

      let link;

      if (this.selectedCount) {
        link = this.req.items[this.selected[0]].url;
      } else {
        link = this.$route.path;
      }

      try {
        const hash = await api.checksum(link, algo);
        event.target.textContent = hash;
      } catch (e) {
        this.$showError(e);
      }
    },
    async calculateDirSize() {
      this.calculatingDirSize = true;
      try {
        let path;
        if (this.selectedCount === 1) {
          path = this.req.items[this.selected[0]].url;
        } else {
          path = this.$route.path;
        }
        
        // Import fetchURL and use it directly for the API call
        const { fetchURL } = await import("@/api/utils");
        const response = await fetchURL(`/api/dirsize${path}`, {});
        const data = await response.json();
        this.dirSize = data.size;
      } catch (e) {
        this.$showError(e);
      } finally {
        this.calculatingDirSize = false;
      }
    },
  },
};
</script>

<style scoped>
.size-bytes {
  color: #666;
  font-size: 0.9em;
  margin-left: 5px;
}

.calculate-size-btn {
  background: var(--blue);
  color: white;
  border: none;
  padding: 4px 12px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.85em;
  margin-left: 10px;
}

.calculate-size-btn:hover {
  background: var(--dark-blue);
}

.calculating-size {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  color: var(--blue);
}

.calculating-size .material-icons {
  font-size: 16px;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.spin {
  animation: spin 1s linear infinite;
}
</style>

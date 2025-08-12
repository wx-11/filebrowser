import { defineStore } from "pinia";

export const useFileStore = defineStore("file", {
  // convert to a function
  state: (): {
    req: Resource | null;
    oldReq: Resource | null;
    reload: boolean;
    selected: number[];
    multiple: boolean;
    isFiles: boolean;
    preselect: string | null;
    scrollPosition: { [path: string]: number };
    calculateDirSizes: boolean;
  } => ({
    req: null,
    oldReq: null,
    reload: false,
    selected: [],
    multiple: false,
    isFiles: false,
    preselect: null,
    scrollPosition: {},
    calculateDirSizes: false,
  }),
  getters: {
    selectedCount: (state) => state.selected.length,
    // route: () => {
    //   const routerStore = useRouterStore();
    //   return routerStore.router.currentRoute;
    // },
    // isFiles: (state) => {
    //   const layoutStore = useLayoutStore();
    //   return !layoutStore.loading && state.route._value.name === "Files";
    // },
    isListing: (state) => {
      return state.isFiles && state?.req?.isDir;
    },
  },
  actions: {
    // no context as first argument, use `this` instead
    toggleMultiple() {
      this.multiple = !this.multiple;
    },
    updateRequest(value: Resource | null) {
      const selectedItems = this.selected.map((i) => this.req?.items[i]);
      this.oldReq = this.req;
      this.req = value;

      this.selected = [];

      if (!this.req?.items) return;
      this.selected = this.req.items
        .filter((item) =>
          selectedItems.some((rItem) => rItem?.url === item.url)
        )
        .map((item) => item.index);
    },
    removeSelected(value: any) {
      const i = this.selected.indexOf(value);
      if (i === -1) return;
      this.selected.splice(i, 1);
    },
    // 保存当前路径的滚动位置
    saveScrollPosition(path: string) {
      this.scrollPosition[path] = window.scrollY;
    },
    // 恢复指定路径的滚动位置
    restoreScrollPosition(path: string) {
      const position = this.scrollPosition[path];
      if (position !== undefined) {
        // 使用 nextTick 确保 DOM 已更新
        setTimeout(() => {
          window.scrollTo(0, position);
        }, 0);
        return true;
      }
      return false;
    },
    // 清除指定路径的滚动位置记录
    clearScrollPosition(path: string) {
      delete this.scrollPosition[path];
    },
    // easily reset state using `$reset`
    clearFile() {
      this.$reset();
    },
  },
});

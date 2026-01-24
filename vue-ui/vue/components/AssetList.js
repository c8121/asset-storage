(function () {
    return {

        parentUrl: '/vue/components/common/BaseObjectWidget.js',

        components: {
            'PathItemTree': '/vue/components/PathItemTree.js',
            'Upload': '/vue/components/Upload.js',
            'Download': '/vue/components/Download.js',
            'MetaData': '/vue/components/MetaData.js',
            'SelectedAssets': '/vue/components/SelectedAssets.js',
        },

        mixins: [
            ui.vueMixin
        ],

        template: `
            <div class="widget AssetList">
                <div v-if="list" class="assets row row-cols-auto g-3 mt-1">
                    <div class="sticky-top bg-white pb-2 pt-2">
                        <div class="row">
                            <div class="col-auto">
                                <div class="pathDropdown">
                                    <div class="input-group input-group-sm">
                                        <button class="btn btn-sm btn-outline-secondary border-light-subtle dropdown-toggle" 
                                            type="button" data-bs-toggle="dropdown" data-bs-auto-close="outside" aria-expanded="false">
                                            {{ pathButtonText }}
                                        </button>
                                        <div class="dropdown-menu" style="overflow: auto;max-height: 90vh;">
                                            <PathItemTree ref="pathItemTree" :value="null"
                                                @click="pathItemClicked"></PathItemTree>
                                        </div>
                                        <button v-if="pathItem" class="btn btn-sm btn-outline-secondary border-light-subtle"
                                            @click="clearPathFilter">x</button>
                                    </div>
                                </div>
                            </div>
                            <div class="col-auto">
                                <div class="input-group input-group-sm">
                                    <label for="page" class="input-group-text border-light-subtle">Page</label>
                                    <input id="page" v-model="page" class="form-control w-auto border-light-subtle text-center" style="width: 50px !important" @changed="pageChanged" @keyup.enter="pageChanged">
                                    <button v-if="page>1" class="btn btn-outline-secondary border-light-subtle" @click="page-=(page>1?1:0);pageChanged()">&lt;</button>
                                    <button v-if="showLoadMore" class="btn btn-outline-secondary border-light-subtle" @click="page++;pageChanged()">&gt;</button>
                                </div>
                            </div>
                            <div class="col-auto">
                                <div class="input-group input-group-sm">
                                    <label for="type" class="input-group-text border-light-subtle">Type</label>
                                    <select v-model="type" class="form-select border-light-subtle" style="max-width: 150px !important" @change="typeChanged">
                                        <option :value="null">Any</option>
                                        <option value="image/*">Image</option>
                                        <option value="image/jpeg">Image: JPG</option>
                                        <option value="image/webp">Image: WEBP</option>
                                        <option value="image/png">Image: PNG</option>
                                        <option value="image/gif">Image: GIF</option>
                                        <option value="audio/*">Audio</option>
                                        <option value="video/*">Video</option>
                                        <option value="application/pdf">PDF</option>
                                        <option :value="null"> - </option>
                                        <option v-for="t in mimeTypes" :value="t.Id.toString()">{{t.Name}}</option>
                                    </select>
                                </div>
                            </div>
                            <div class="col-auto">
                                <div class="input-group input-group-sm">
                                    <input id="findName" v-model="findName" class="form-control w-auto border-light-subtle border-end-0"
                                            @keyup.enter="findByName">
                                    <button v-if="findName" class="btn btn-sm btn-outline-secondary border-light-subtle border-start-0"
                                            @click="clearFindFilter">x</button>
                                    <button class="btn btn-outline-secondary border-light-subtle" 
                                            @click="findByName">find...</button>
                                </div>
                            </div>
                            <div class="col-auto" v-if="Object.keys(selectedAssets).length > 0">
                                <button class="btn btn-sm btn-outline-secondary border-light-subtle"
                                    @click="$refs.selectionToast.showToast()">{{ Object.keys(selectedAssets).length }} files selected</button>
                            </div>
                            <div class="col-auto">
                                <button class="btn btn-sm btn-outline-secondary border-light-subtle"
                                    @click="$refs.uploadToast.showToast()">Upload</button>
                            </div>
                        </div>
                    </div>
                    <template v-for="asset in list">
                        <div v-if="asset.nextGroup" class="asset-group-header">{{ asset.groupKey }}</div>
                        <div class="asset col">
                            <div class="card bg-light">
                                <div class="card-body">
                                    <div class="asset-filename text-center p-3 small">
                                        {{ asset.Name }}
                                    </div>
                                    <div class="asset-image">
                                        <img @click="showDownload(asset)"
                                            role="button"
                                            class="card-img-top asset-preview not-ready" 
                                            :src="'/assets/thumbnail/' + asset.Hash"
                                            :alt="asset.Name" :title="asset.Name" />
                                    </div>
                                </div>
                                <div class="card-footer text-end p-0">
                                    <input type="checkbox" :checked="isSelected(asset)" @change="selectAssetClick(asset, $event)">
                                    <button @click="showMetaData(asset)" class="btn btn-sm"><i>M</i></button>
                                </div>
                            </div>
                        </div>
                    </template>
                </div>
                <div class="row m-2">
                    <div class="col text-start"></div>
                    <div class="col text-center">
                        <button v-if="showLoadMore" @click="loadMore" class="btn btn-light" id="loadMore">Load more...</button>
                    </div>
                    <div class="col text-end"></div>
                </div>


                <div class="position-fixed toast-container top-0 end-0 p-3">
                    <Download ref="downloadToast" :value="selectedAssetMetaData" 
                        @file-click="downloadAsset"></Download>
                    <Upload ref="uploadToast" 
                        @meta-data-click="showMetaData"></Upload>
                    <MetaData ref="metaDataToast" :value="selectedAssetMetaData" 
                        @file-click="downloadAsset(selectedAssetMetaData)"></MetaData>
                    <SelectedAssets ref="selectionToast" :value="selectedAssets" 
                        @file-click="downloadAsset"
                        @remove-click="unselectAsset"
                        @meta-data-click="showMetaData"></SelectedAssets>
                </div>

            </div>
            `,

        data() {
            return {
                list: [],

                offset: 0,
                count: 30,
                page: 1,
                type: null,
                pathItem: null,
                findName: null,

                mimeTypes: [],

                loading: false,
                showLoadMore: true,

                pathButtonText: 'Paths',

                selectedAssetMetaData: null,

                selectedAssets: {}
            }
        },
        methods: {
            loadDataObject() {
                const self = this;
                self.loadAssetList();
                self.loadMimeTypes();
            },

            loadAssetList() {
                const self = this;
                self.loading = true;

                client.post('/assets/list', self.getListFilter()).then((json) => {
                    if (json) {
                        self.list = json;
                        self.showLoadMore = json.length >= self.count
                        self.page = (self.offset + self.count) / self.count;
                        self.createGroupKeys(self.list);
                    }
                    self.loading = false;
                    self.enableAutoLoadMore();
                });
            },

            reloadAssetList(offset) {
                const self = this;
                self.list.splice(0, self.list.length);
                self.offset = offset;
                self.loadAssetList();
            },

            loadMore() {
                const self = this;
                self.loading = true;

                self.offset += self.count
                const filter = self.getListFilter();

                client.post('/assets/list', filter).then((json) => {
                    if (json) {
                        self.showLoadMore = json.length > 0
                        self.page = (self.offset + self.count) / self.count;
                        self.createGroupKeys(json);
                        for (const item of json)
                            self.list.push(item)
                    }
                    self.loading = false;
                });
            },

            pageChanged() {
                const self = this;
                self.reloadAssetList((self.page - 1) * self.count);
            },

            typeChanged() {
                const self = this;
                self.reloadAssetList(0);
            },

            findByName() {
                const self = this;
                self.reloadAssetList(0);
            },

            clearFindFilter() {
                const self = this;
                self.findName = null;
                self.reloadAssetList(0);
            },

            getListFilter() {
                const self = this;

                return {
                    Offset: self.offset,
                    Count: self.count,
                    MimeType: self.type,
                    FileName: self.findName,
                    //PathName: null,
                    PathId: self.pathItem
                };
            },

            createGroupKeys(list) {
                let lastGroup = "";
                for (const asset of list) {
                    if (!asset.groupKey) {
                        asset.groupKey = asset.FileTime.substring(0, 10);
                        if (asset.groupKey !== lastGroup) {
                            asset.nextGroup = true;
                            lastGroup = asset.groupKey;
                        } else {
                            asset.nextGroup = false;
                        }
                    }
                }
            },

            pathItemClicked(item) {
                const self = this;
                self.pathItem = item.Id;
                self.reloadAssetList(0);
                self.pathButtonText = item.Name;
            },

            clearPathFilter() {
                const self = this;
                self.pathItem = null;
                self.reloadAssetList(0);
                self.pathButtonText = 'Paths';
            },

            selectAssetClick(asset, e) {
                const self = this;
                if (e.target.checked) {
                    self.selectedAssets[asset.Hash] = asset;
                } else {
                    delete self.selectedAssets[asset.Hash];
                }
                if (Object.keys(self.selectedAssets).length > 0)
                    self.$refs.selectionToast.showToast();
                else
                    self.$refs.selectionToast.hideToast();
            },

            unselectAsset(asset) {
                const self = this;
                delete self.selectedAssets[asset.Hash];
                if (Object.keys(self.selectedAssets).length > 0)
                    self.$refs.selectionToast.showToast();
                else
                    self.$refs.selectionToast.hideToast();
            },

            isSelected(asset) {
                const self = this;
                return self.selectedAssets.hasOwnProperty(asset.Hash);
            },

            downloadAsset(asset) {
                window.open('/assets/' + asset.Hash);
            },

            showDownload(asset) {
                const self = this;
                client.get('/assets/metadata/' + asset.Hash).then((json) => {
                    self.selectedAssetMetaData = json;
                    self.$refs.downloadToast.showToast();
                });
            },

            showMetaData(asset) {
                const self = this;
                client.get('/assets/metadata/' + asset.Hash).then((json) => {
                    self.selectedAssetMetaData = json;
                    self.$refs.metaDataToast.showToast();
                });
            },

            //Autoload more items when scrolling to bottom
            enableAutoLoadMore() {
                const self = this;

                if (!self.scrollListener) {

                    self.scrollListener = (e) => {

                        if (self.loading || !self.showLoadMore)
                            return;

                        const loadMoreButton = document.getElementById("loadMore");
                        const rect = loadMoreButton.getBoundingClientRect();
                        const elemTop = rect.top;
                        const elemBottom = rect.bottom;

                        // Only completely visible elements return true:
                        const isVisible = (elemTop > 0) && (elemBottom <= window.innerHeight);
                        // Partially visible elements return true:
                        //const isVisible = elemTop < window.innerHeight && elemBottom >= 0;
                        if (isVisible) {
                            self.loadMore();
                        }
                    }
                    document.addEventListener('scroll', self.scrollListener);
                }
            },

            loadMimeTypes() {
                const self = this;

                client.get("/mimetypes/list").then((json) => {
                    if (json) {
                        self.mimeTypes.splice(0, self.mimeTypes.length);
                        for (const item of json) {
                            self.mimeTypes.push(item);
                        }
                    }
                });

            }
        }
    }
})();
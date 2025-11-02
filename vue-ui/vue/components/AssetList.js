(function () {
    return {

        parentUrl: '/vue/components/common/BaseObjectWidget.js',

        mixins: [
            ui.vueMixin
        ],

        template: `
            <div class="widget AssetList">
                <div v-if="list" class="assets row row-cols-auto g-3 mt-1">
                    <div class="sticky-top bg-white pb-2 pt-2">
                        <div class="row">
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
                                    <input id="findName" v-model="findName" class="form-control w-auto border-light-subtle" @keyup.enter="findByName">
                                    <button class="btn btn-outline-secondary border-light-subtle" @click="findByName">find...</button>
                                </div>
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
                                        <img @click="selectAsset(asset)"
                                            role="button"
                                            class="card-img-top asset-preview not-ready" 
                                            :src="'/assets/thumbnail/' + asset.Hash"
                                            :alt="asset.Name" :title="asset.Name" />
                                    </div>
                                </div>
                                <div class="card-footer text-end p-0">
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
                showLoadMore: true
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

            loadMore() {
                const self = this;
                self.loading = true;

                self.offset += self.count

                client.post('/assets/list', self.getListFilter()).then((json) => {
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

                self.list.splice(0, self.list.length);
                self.offset = (self.page - 1) * self.count
                self.loadAssetList()
            },

            typeChanged() {
                const self = this;

                self.list.splice(0, self.list.length);
                self.offset = 0;
                self.loadAssetList()
            },

            findByName() {
                const self = this;

                self.list.splice(0, self.list.length);
                self.offset = 0;
                self.loadAssetList()
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

            setPathItemId(id) {
                const self = this;
                self.pathItem = id;
                self.list.splice(0, self.list.length);
                self.offset = 0;
                self.loadAssetList()
            },

            selectAsset(asset) {
                this.downloadAsset(asset);
            },

            downloadAsset(asset) {
                window.open('/assets/' + asset.Hash);
            },

            showMetaData(asset) {
                window.open('/assets/metadata/' + asset.Hash);
            },

            //Autoload more items when scrolling to bottom
            enableAutoLoadMore() {
                const self = this;

                const loadMoreButton = document.getElementById("loadMore");
                self.scrollListener = (e) => {
                    if (self.loading || !self.showLoadMore)
                        return;

                    const rect = loadMoreButton.getBoundingClientRect();
                    const elemTop = rect.top;
                    const elemBottom = rect.bottom;

                    // Only completely visible elements return true:
                    var isVisible = (elemTop >= 0) && (elemBottom <= window.innerHeight);
                    // Partially visible elements return true:
                    //isVisible = elemTop < window.innerHeight && elemBottom >= 0;
                    if (isVisible)
                        self.loadMore();
                }
                document.addEventListener('scroll', self.scrollListener);
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
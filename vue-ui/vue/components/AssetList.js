(function () {
    return {

        parentUrl: '/vue/components/common/BaseObjectWidget.js',
        
        mixins: [
            ui.vueMixin  
        ],

        template: `
            <div class="widget AssetList">
                <div v-if="list" class="assets row row-cols-auto g-3 mt-1">
                    <template v-for="asset in list">
                        <div v-if="asset.nextGroup" class="asset-group-header">{{ asset.groupKey }}</div>
                        <div class="asset col">
                            <div class="card bg-light">
                                <div class="card-body">
                                    <div class="asset-filename text-center p-3 small">
                                        {{ asset.Name }}
                                    </div>
                                    <div class="asset-image">
                                        <img  @click="selectAsset(asset)"
                                            role="button"
                                            class="card-img-top asset-preview not-ready" 
                                            :src="'/assets/thumbnail/' + asset.Hash" />
                                    </div>
                                </div>
                                <div class="card-footer">
                                    
                                </div>
                            </div>
                        </div>
                    </template>
                </div>
                <div class="row m-2">
                    <div class="col text-start"></div>
                    <div class="col text-center">
                        <button @click="loadMore" class="btn btn-light" id="loadMore">Load more...</button>
                    </div>
                    <div class="col text-end"></div>
                </div>
            </div>
            `,

        data() {
            return {
                REST_URL: '/assets/list',
                list: [],

                offset: 0,
                count: 30,

                loading: false
            }
        },
        methods: {
            loadDataObject() {
                const self = this;
                self.loading = true;

                client.get(self.REST_URL + '/' + self.offset + '/' + self.count).then((json) => {
                    if(json) {
                        self.list = json;
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

                client.get(self.REST_URL + '/' + self.offset + '/' + self.count).then((json) => {
                    if(json) {
                        self.createGroupKeys(json);
                        for(const item of json)
                            self.list.push(item)
                    }
                    self.loading = false;
                });
            },

            createGroupKeys(list) {
                let lastGroup = "";
                for (const asset of list) {
                    if (!asset.groupKey) {
                        asset.groupKey = asset.FileTime.substring(0, 10);
                        if (asset.groupKey != lastGroup) {
                            asset.nextGroup = true;
                            lastGroup = asset.groupKey;
                        } else {
                            asset.nextGroup = false;
                        }
                    }
                }
            },

            selectAsset(asset) {
                this.downloadAsset(asset);
            },

            downloadAsset(asset) {
                window.open('/assets/' + asset.Hash);
            },

            //Autoload more items when scrolling to bottom
            enableAutoLoadMore() {
                const self = this;

                const loadMoreButton = document.getElementById("loadMore");
                self.scrollListener = (e) => {
                    if (self.loading)
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
            }
        }
    }
})();
(function () {
    return {

        parentUrl: '/vue/components/common/BaseObjectWidget.js',

        components: {},

        mixins: [
            ui.vueMixin
        ],

        template: `
            <div class="widget CollectionList">
                <div v-if="list" class="collections">
                                        
                    <template v-for="collection in list">
                        <div class="row row-cols-auto g-3 mt-1 mb-3 bg-light border-bottom">
                            <div class="col-6 text-primary">{{collection.Name}}</div>
                            <div class="col">{{formatter.date(collection.Created)}}</div>
                            <div class="col">{{collection.Assets.length}} assets</div>
                        </div>
                        <div v-if="collection.Description" class="row">
                            <div class="col ps-5 pt-3 text-secondary"><pre>{{collection.Description}}</pre></div>
                        </div>
                        <div class="assets">
                            <template v-for="asset in collection.Assets">
                                <div class="asset col">
                                    <div class="card bg-light">
                                        <div class="card-body">
                                            <div class="asset-image">
                                                <img @click="showDownload(asset)"
                                                    role="button"
                                                    class="card-img-top asset-preview not-ready" 
                                                    :src="'/assets/thumbnail/' + asset"
                                                     />
                                            </div>
                                        </div>
                                        <div class="card-footer text-end p-0">
                                            <button @click="showMetaData(asset)" class="btn btn-sm"><i>M</i></button>
                                        </div>
                                    </div>
                                </div>
                            </template>
                        </div>
                    </template>
                    
                    <pre>{{ list }}</pre>
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

                loading: false,
                showLoadMore: true,

                formatter: Formatter
            }
        },
        methods: {
            loadDataObject() {
                const self = this;
                self.loadCollectionList();
            },

            loadCollectionList() {
                const self = this;
                self.loading = true;

                client.post('/collections/list', self.getListFilter()).then((json) => {
                    if (json) {
                        self.list.splice(0, self.list.length);
                        for (const item of json) {
                            self.loadAssets(item).then((asset) => {
                                item.Description = asset.Description;
                                item.Assets = asset.Assets;
                                self.list.push(item);
                            });
                        }
                        self.showLoadMore = json.length >= self.count
                    }
                    self.loading = false;
                    self.enableAutoLoadMore();
                });
            },

            loadMore() {
                const self = this;
                self.loading = true;

                self.offset += self.count
                const filter = self.getListFilter();

                client.post('/collections/list', filter).then((json) => {
                    if (json) {
                        self.showLoadMore = json.length > 0
                        for (const item of json) {
                            self.loadAssets(item).then((asset) => {
                                item.Description = asset.Description;
                                item.Assets = asset.Assets;
                                self.list.push(item);
                            });
                        }
                    }
                    self.loading = false;
                });
            },

            loadAssets(collection) {
                const self = this;

                return client.get('/collections/' + collection.Hash).then((json) => {
                    return json;
                });
            },

            getListFilter() {
                const self = this;

                return {
                    Offset: self.offset,
                    Count: self.count,
                };
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
        }
    }
})();
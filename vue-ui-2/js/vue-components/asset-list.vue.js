export default {
    template: `
        <div>
        
            <div class="row asset-list">
                <div v-for="asset in list" class="asset col pb-3">
                    <div class="card bg-light">
                        <div class="card-body">
                            <div class="text-center p-1 small">
                                {{ asset.Name }}
                            </div>
                            <div>
                                <img @click="assetClick(asset)"
                                    role="button"
                                    class="card-img-top asset-preview not-ready" 
                                    :src="'/assets/thumbnail/' + asset.Hash"
                                    :alt="asset.Name" :title="asset.Name" />
                            </div>
                        </div>
                        <div class="card-footer">
                            <div class="row">
                                <div class="col">
                                    <div class="form-check form-switch">
                                        <input type="checkbox"  class="form-check-input" role="switch" @change="assetSelect(asset, $event)">
                                    </div>
                                </div>
                                <div class="col text-end">
                                    <button class="btn btn-sm btn-outline-secondary" @click="downloadClick(asset)"><i class="bi bi-download"></i></button>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
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

    props: {
        filter: {
            type: Object,
            default: null
        }
    },

    data() {
        return {
            list: [],

            offset: 0,
            count: 30,

            showLoadMore: true
        }
    },

    methods: {

        loadAssets() {

            const self = this;
            self.offset = 0;
            fetch('/assets/list', self.createRequestOptions())
                .then(res => res.json())
                .then(json => {
                    self.list = json;
                });
        },

        loadMore() {
            const self = this;

            self.offset += self.count
            fetch('/assets/list', self.createRequestOptions())
                .then(res => res.json())
                .then(json => {
                    for(const item of json) {
                        self.list.push(item);
                    }
                });
        },

        createRequestOptions() {
            const self = this;

            const listFilter = {
                Offset: self.offset,
                Count: self.count,
                MimeType: null,
                //FileName: null,
                //PathName: null,
                //PathId: null,
                //Face: null
            }

            for(const p of Object.keys(self.filter)) {
                const value = self.filter[p];
                if(value || value === 0)
                    listFilter[p] = value;
            }

            return {
                method: 'POST',
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify(listFilter)
            };
        },

        assetSelect(asset, e) {
            const self = this;
            if (e.target.checked) {
                this.$emit('componentEvent', 'assetSelect', 'asset-list', asset);
            } else {
                this.$emit('componentEvent', 'assetUnselect', 'asset-list', asset);
            }
        },

        assetClick(asset) {
            this.$emit('componentEvent', 'assetClick', 'asset-list', asset);
        },

        downloadClick(asset) {
            this.$emit('componentEvent', 'assetDownloadClick', 'asset-list', asset);
        }
    },

    emits: ['componentEvent'],

    created() {
        this.loadAssets();
    }
}
export default {
    template: `
        <div>
            <div v-if="value">
                <div @click="assetClick()" role="button">
                    <span class="text-primary"><strong>{{ value.Name }}</strong></span>
                    <span class="text-secondary ps-3">{{ value.MimeType }}</span>
                    <span class="text-secondary ps-3">{{ value.FileTime }}</span>
                </div>
            </div>
            <div v-if="asset && asset.Origins">
                <div v-for="origin in asset.Origins">
                    <div class="text-primary mt-3 text-truncate">{{ origin.Path }}</div>
                    <div class="text-secondary">{{ origin.Owner }}</div>
                </div>
            </div>
            
            <div class="mt-3">
                <div class="input-group">
                    <select class="form-select" v-model="filterName" @change="filterChange">
                        <option :value="null">Original</option>
                        <option value="image">Image</option>
                    </select>
                    <button class="btn btn-primary"
                        @click="downloadClick()">
                       <i class="bi bi-download"></i>
                    </button>
                </div>
            </div>
            <div class="mt-3">
                <form method="post" ref="filterRequest" target="_blank">
                    <div v-if="filterParams" v-for="p in filterParams">
                        <div class="input-group input-group-sm mb-1">
                            <label class="input-group-text">{{p.label}}</label>
                            <input type="number" :name="p.name" v-model="p.value" class="form-control">
                        </div>
                    </div>
                </form>
            </div>
        </div>
    `,

    props: {
        value: {
            type: Object,
            default: null
        }
    },

    data() {
        return {
            asset: null,

            filterName: null,
            filterParams: null,

            availableFilterParams: {
                image: [
                    { name: 'width', label: 'Width', value: 100 },
                    { name: 'height', label: 'Height', value: "" }
                ]
            }
        }
    },

    watch: {
        value() {
            this.loadAsset();
        }
    },

    methods: {

        loadAsset() {
            const self = this;
            if(!self.value || !self.value.Hash)
                return;

            const requestOptions = {
                method: 'GET'
            }
            fetch('/assets/metadata/' + self.value.Hash, requestOptions)
                .then(res => res.json())
                .then(json => {
                    self.asset = json;
                });
        },

        filterChange() {
            const self = this;
            self.filterParams = null;
            if(!self.filterName || !self.availableFilterParams[self.filterName]) {
                self.filterParams = self.filterName;
            }
            self.filterParams = self.availableFilterParams[self.filterName];
        },

        assetClick() {
            this.$emit('componentEvent', 'assetClick', 'asset-view', this.value);
        },

        downloadClick() {
            const self = this;
            if(!self.filterName) {
                this.$emit('componentEvent', 'assetDownloadClick', 'asset-view', this.value);
            } else {
                self.$refs.filterRequest.action =
                    "/assets/filter/" + self.filterName +
                    "/" + self.value.Hash;
                self.$refs.filterRequest.submit();
            }
        }
    },

    emits: ['componentEvent'],
}
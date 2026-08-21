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
                        <div class="card-footer text-end">
                            <div class="form-check form-switch">
                                <input type="checkbox"  class="form-check-input" role="switch" @change="assetSelect(asset, $event)">
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    `,

    data() {
        return {
            list: [],
        }
    },

    methods: {

        loadAssets() {
            const self = this;

            const requestOptions = {
                method: 'POST',
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify({
                    Offset: 0,
                    Count: 30,
                    MimeType: null,
                    FileName: null,
                    //PathName: null,
                    PathId: null,
                    Face: null
                })
            }
            fetch('/assets/list', requestOptions)
                .then(res => res.json())
                .then(json => {
                    self.list = json;
                });
        },

        assetSelect(asset, e) {
            const self = this;
            if (e.target.checked) {
                this.$emit('componentEvent', 'assetSelect', asset);
            } else {
                this.$emit('componentEvent', 'assetUnselect', asset);
            }
        },

        assetClick(asset) {
            this.$emit('componentEvent', 'assetClick', 'asset-list', asset);
        }
    },

    emits: ['componentEvent'],

    created() {
        this.loadAssets();
    }
}
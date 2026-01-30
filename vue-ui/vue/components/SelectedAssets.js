(function () {
    return {

        template: `
            <div :class="'toast ' + showToastCss" role="alert" aria-live="assertive" aria-atomic="true">
                <div class="toast-header">
                    <strong class="me-auto">{{ headerCaption }}</strong>
                    <button type="button" class="btn-close" aria-label="Close" @click="hideToast"></button>
                </div>
                <div class="toast-body">
                    <div v-if="value" style="overflow: auto;max-height: 50vh;">
                        
                        <div v-for="(asset, hash) in value" class="position-relative">
                            <p class="text-primary" role="button"
                                @click="onFileClick(asset)">{{ asset.Name }}</p>
                            <div class="position-absolute top-0 end-0">
                                <button class="btn btn-sm btn-light"
                                    @click="onMetaDataClick(asset)">M</button>
                                <button class="btn btn-sm btn-light"
                                    @click="onRemoveClick(asset)">x</button>
                            </div>
                        </div>

                        <!-- pre>{{ value }}</pre -->
                    </div>
                    <div v-if="value">
                        <button class="btn btn-sm btn-link"
                            @click="collectionParamsVisible=!collectionParamsVisible"
                            >{{ showCreateCollectionButtonCaption }}</button>

                        <div v-if="collectionParamsVisible">
                            <div class="input-group input-group-sm mt-1">
                                <label class="input-group-text">Name</label>
                                <input v-model="collectionName" class="form-control">
                            </div>
                            <div class="input-group input-group-sm mt-1">
                                <label class="input-group-text">Description</label>
                                <textarea v-model="collectionDescription" class="form-control"></textarea>
                            </div>
                            <button class="btn btn-sm btn-primary mt-1"
                                @click="createCollection"
                                >{{ createCollectionButtonCaption }}</button>
                        </div>
                    </div>
                </div>
            </div>`,

        props: {
            headerCaption: {
                type: String,
                default: "Selected Files"
            },
            showCreateCollectionButtonCaption: {
                type: String,
                default: "Create Collection"
            },
            createCollectionButtonCaption: {
                type: String,
                default: "Create Collection"
            },
            value: {
                type: Object,
                default: {}
            }
        },

        data() {
            return {
                showToastCss: '',

                collectionParamsVisible: false,
                collectionName: '',
                collectionDescription: '',
            }
        },
        
        methods: {
            showToast() {
                this.showToastCss = 'show';
            },
            hideToast() {
                this.showToastCss = '';
            },
            onRemoveClick(asset) {
                this.$emit('removeClick', asset);    
            },
            onFileClick(asset) {
                this.$emit('fileClick', asset);
            },
            onMetaDataClick(asset) {
                this.$emit('metaDataClick', asset);
            },
            createCollection() {
                const self = this;
                const query = {
                    Name: self.collectionName,
                    Description: self.collectionDescription,
                    AssetHashes: []
                }
                for(const hash in self.value)
                    query.AssetHashes.push(hash);
                
                client.post('/collections/add', query).then((json) => {
                    console.log(json);
                });
            }
        },
        emits: [
            'removeClick',
            'fileClick',
            'metaDataClick'
        ]
    }
})();
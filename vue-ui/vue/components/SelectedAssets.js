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

                        <!-- <pre>{{ value }}</pre> -->
                    </div>
                </div>
            </div>`,

        props: {
            headerCaption: {
                type: String,
                default: "Selected Files"
            },
            value: {
                type: Object,
                default: {}
            }
        },

        data() {
            return {
                showToastCss: ''
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
            }
        },
        emits: [
            'removeClick',
            'fileClick',
            'metaDataClick'
        ]
    }
})();
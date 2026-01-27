(function () {
    return {

        template: `
            <div :class="'toast ' + showToastCss" role="alert" aria-live="assertive" aria-atomic="true">
                <div class="toast-header">
                    <strong class="me-auto">{{ headerCaption }}</strong>
                    <button type="button" class="btn-close" aria-label="Close" @click="hideToast"></button>
                </div>
                <div class="toast-body">
                    <div v-if="value">

                        <div>
                            <span class="text-primary">{{ value.Origins[0].Name }}</span>
                            <span class="text-secondary ps-3 text-nowrap">{{ value.MimeType }}</span>
                        </div>

                        <div class="mt-3">
                            <div class="input-group">
                                <select class="form-select" v-model="filterName" @change="filterChange">
                                    <option :value="null">Original</option>
                                    <option value="image">Image</option>
                                </select>
                                <button class="btn btn-primary"
                                    @click="onDownloadClick()">
                                    {{ buttonCaption }}
                                </button>
                            </div>
                        </div>
                        <div class="mt-3">
                            <form method="post" ref="filterRequest" target="_blank">
                                <div v-if="filterParams" v-for="p in filterParams">
                                    <div class="input-group mb-1">
                                        <label class="input-group-text">{{p.label}}</label>
                                        <input type="number" :name="p.name" :value="p.value" class="form-control">
                                    </div>
                                </div>
                            </form>
                        </div>

                        <!-- <pre>{{ value }}</pre> -->
                    </div>
                </div>
            </div>`,

        props: {
            headerCaption: {
                type: String,
                default: "Download"
            },
            buttonCaption: {
                type: String,
                default: "Download"
            },
            value: {
                type: Object,
                default: null
            }
        },

        data() {
            return {
                showToastCss: '',

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
        methods: {
            showToast() {
                this.showToastCss = 'show';
            },
            hideToast() {
                this.showToastCss = '';
            },
            filterChange() {
                const self = this;
                self.filterParams = null;
                if(!self.filterName || !self.availableFilterParams[self.filterName]) {
                    self.filterParams = self.filterName;
                }
                self.filterParams = self.availableFilterParams[self.filterName];
            },
            onDownloadClick() {
                const self = this;
                if(!self.filterName) {
                    this.$emit('fileClick', this.value);    
                } else {
                    self.$refs.filterRequest.action =
                        "/assets/filter/" + self.filterName +
                        "/" + self.value.Hash;
                    self.$refs.filterRequest.submit();
                }
            }
        },
        emits: [
            'fileClick'
        ]
    }
})();
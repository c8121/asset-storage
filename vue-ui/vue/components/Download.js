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

                        <p class="small text-secondary">{{ value.Hash }}</p>
                        <p class="small">{{ value.MimeType }}</p>

                        <p>
                            <button class="btn btn-primary"
                                @click="onFileClick()">
                                {{ buttonCaption }}
                            </button>
                        </p>

                        <p><a class="small" 
                            data-bs-toggle="collapse"
                            href="#downloadFiltered" role="button">Download filtered (not availabe yet)</a></p>
                        <div class="collapse" id="downloadFiltered">
                            <form method="post" ref="filterRequest" target="_blank">
                                <div class="mt-5">
                                    <select class="form-select" v-model="filterName">
                                        <option :value="null"></option>
                                        <option value="image">Image</option>
                                    </select>
                                </div>
                                <div class="row mt-1">
                                    <div class="col-auto">Params</div>
                                    <div class="col-auto">
                                        <textarea name="params" v-model="filterParams" class="form-control"></textarea>
                                    </div>
                                </div>
                            </form>
                            <div class="mt-1">
                                <button class="btn btn-primary"
                                    @click="onFilterClick()">
                                    {{ filterButtonCaption }}
                                </button>
                            </div>
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
                default: "Download original"
            },
            filterButtonCaption: {
                type: String,
                default: "Download filtered"
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
                filterParams: ""
            }
        },
        methods: {
            showToast() {
                this.showToastCss = 'show';
            },
            hideToast() {
                this.showToastCss = '';
            },
            onFileClick() {
                this.$emit('fileClick', this.value);
            },
            onFilterClick() {
                const self = this;
                self.$refs.filterRequest.action = 
                    "/assets/filter/" + self.filterName + 
                    "/" + self.value.Hash;
                self.$refs.filterRequest.submit();
            }
        },
        emits: [
            'fileClick'
        ]
    }
})();
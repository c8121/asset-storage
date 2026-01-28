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
                        <div @click="onFileClick()" role="button">
                            <span class="text-primary">{{ value.Origins[0].Name }}</span>
                            <span class="text-secondary ps-3 text-nowrap">{{ value.MimeType }}</span>
                        </div>

                        <div class="mt-3">
                            <div v-for="origin in value.Origins" class="border-top border-light-subtle p-1 mb-1">
                                <a
                                        @click="onFileClick()"
                                        role="button">
                                    {{ formatter.date(origin.FileTime) }}<br />
                                    <small class="text-secondary pe-1">{{ origin.Path }}</small>
                                    <strong class="text-primary">{{ origin.Name }}</strong>
                                </a>
                            </div>
                        </div>
                        
                        <div class="mt-3 small text-secondary">{{ value.Hash }}</div>

                        <!-- <pre>{{ value }}</pre> -->
                    </div>
                </div>
            </div>`,

        props: {
            headerCaption: {
                type: String,
                default: "Meta-Data"
            },
            value: {
                type: Object,
                default: null
            }
        },

        data() {
            return {
                showToastCss: '',

                formatter: Formatter
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
            }
        },
        emits: [
            'fileClick'
        ]
    }
})();
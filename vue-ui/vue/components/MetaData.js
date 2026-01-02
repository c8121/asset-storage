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

                        <div v-for="origin in value.Origins">
                            <p class="text-primary"
                                    @click="onFileClick()"
                                    role="button">
                                {{ origin.FileTime }}<br />
                                <strong>{{ origin.Name }}</strong><br />
                                <small>{{ origin.Path }}</small>
                            </p>
                        </div>

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
            onFileClick() {
                this.$emit('fileClick', this.value);
            }
        },
        emits: [
            'fileClick'
        ]
    }
})();
(function () {
    return {

        template: `
            <div :class="'toast ' + showToastCss" role="alert" aria-live="assertive" aria-atomic="true">
                <div class="toast-header">
                    <strong class="me-auto">{{ headerCaption }}</strong>
                    <button type="button" class="btn-close" aria-label="Close" @click="hideToast"></button>
                </div>
                <div class="toast-body">
                    <div v-if="faces" class="d-flex flex-wrap justify-content-center">
                        <div class="p-1" v-for="face in faces">
                            <img :src="'/faces/' + face"
                                class="rounded-circle"
                                style="min-width: 75px; max-width: 75px">
                        </div>
                    </div>
                </div>
            </div>`,

        props: {
            headerCaption: {
                type: String,
                default: "Faces"
            },
            value: {
                type: Object,
                default: null
            }
        },

        data() {
            return {
                showToastCss: '',
                faces:[]
            }
        },
        methods: {
            showToast() {
                this.showToastCss = 'show';
            },
            hideToast() {
                this.showToastCss = '';
            },
            loadFaces(asset) {
                const self = this;
                client.get('/faces/' + asset.Hash).then((json) => {
                    self.faces = json
                });
            }
        },
        emits: [],

        created() {
            const self = this;
            self.$watch(
                () => this.value,
                (value) => {
                    self.loadFaces(value);
                }
            )
        },
    }
})();
export default {
    template: `
        <div>
            <div class="overflow-auto" style="max-height: 50vh">
                <div class="row" v-for="(item, hash) in value" :key="hash">
                    <div class="col pt-1 pb-1 border-bottom">{{ item.Name }}</div>
                </div>
            </div>
            <div class="mt-3">
                <button class="btn btn-secondary" @click="createCollectionClick">Create collection</button>
            </div>
        </div>
    `,

    props: {
        value: {
            type: Object,
            default: null
        }
    },

    methods: {
        createCollectionClick() {
            this.$emit('componentEvent', 'createCollectionClick', 'asset-selection-list');
        }
    },

    emits: ['componentEvent'],
}
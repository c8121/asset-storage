export default {
    template: `
        <div>
            <div v-if="value">
                <div @click="assetClick(value)" role="button">
                    <span class="text-primary">{{ value.Name }}</span>
                    <span class="text-secondary ps-3 text-nowrap">{{ value.MimeType }}</span>
                </div>
            </div>
            <pre>{{ value }}</pre>
        </div>
    `,

    props: {
        value: {
            type: Object,
            default: null
        }
    },

    methods: {
        assetClick(asset) {
            this.$emit('componentEvent', 'assetClick', 'asset-view', asset);
        }
    },

    emits: ['componentEvent'],
}
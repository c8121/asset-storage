export default {
    template: `
        <div>
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
    },

    emits: ['componentEvent'],
}
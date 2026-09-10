export default {
    template: `
        <div>
        
            <div class="row collection-list">
                <div v-for="collection in list" class="asset col pb-3">
                    <pre @click="collectionClick(collection)">{{ collection }}</pre>
                </div>
            </div>
        </div>
    `,

    data() {
        return {
            list: [],

            offset: 0,
            count: 30,

            loading: false
        }
    },

    methods: {

        loadCollections() {

            const self = this;
            self.loading = true;

            self.offset = 0;
            fetch('/collections/list', self.createRequestOptions())
                .then(res => res.json())
                .then(json => {
                    self.list = json;
                    self.loading = false;
                });
        },

        createRequestOptions() {
            const self = this;

            const listFilter = {
                Offset: self.offset,
                Count: self.count
            }

            return {
                method: 'POST',
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify(listFilter)
            };
        },

        collectionClick(collection) {
            this.$emit('componentEvent', 'collectionClick', 'collection-list', collection);
        }
    },

    emits: ['componentEvent'],

    created() {
        this.loadCollections();
    }
}
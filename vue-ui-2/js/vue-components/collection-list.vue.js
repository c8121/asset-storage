export default {
    template: `
        <div>
            <div class="collection-list overflow-auto" style="max-height: 50vh" >
                <div v-for="collection in list" class="row pt-3 border-bottom" @click="collectionClick(collection)" role="button">
                    <div class="col">{{ collection.Name }}</div>
                    <div class="col">{{ collection.Created }}</div>
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
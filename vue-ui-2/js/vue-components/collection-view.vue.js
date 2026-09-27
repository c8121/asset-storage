export default {
    template: `
        <div>
            <div v-if="collection">
                <div class="input-group">
                    <label class="input-group-text">Name</label>
                    <input class="form-control" v-model="collection.Name" v-on:keyup="changed = true">
                </div>
                <div class="small p-2">
                    {{ collection.Assets.length }} Assets, {{ formatDate(collection.Created) }}<br>
                </div>
                <div>
                    <label>Description</label>
                    <textarea v-model="collection.Description" class="form-control"  v-on:keyup="changed = true"></textarea>
                </div>
                <div v-if="changed" class="mt-3 text-end">
                    <button class="btn btn-secondary" @click="saveChanges">Save</button>
                </div>
                
            </div>
        </div>
    `,

    props: {
        value: {
            type: Object,
            default: null
        }
    },

    data() {
        return {
            collection: null,

            changed: false,
            dateFormatter: Intl.DateTimeFormat('de-DE')
        }
    },

    watch: {
        value() {
            this.loadCollection();
        }
    },

    methods: {

        loadCollection() {
            const self = this;
            if(!self.value || !self.value.UUID)
                return;

            const requestOptions = {
                method: 'GET'
            }
            fetch('/collections/' + self.value.UUID, requestOptions)
                .then(res => res.json())
                .then(json => {
                    self.collection = json;
                    self.changed = false;
                });
        },

        saveChanges() {
            const self = this;

            const requestParams = {
                method: 'POST',
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify(self.collection)
            }

            fetch('/collections/add', requestParams)
                .then(res => res.json())
                .then(json => {
                    self.collection = json;
                    self.changed = false;

                    this.$emit('componentEvent', 'collectionSaved', 'collection-view', self.collection);
                });
        },

        formatDate(d) {
            return this.dateFormatter.format(Date.parse(d));
        }
    },

    emits: ['componentEvent'],
}
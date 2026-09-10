export default {
    template: `
        <div>
            <pre>{{ value }}</pre>
            <div>
                <button @click="createCollection">Create collection</button>
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
        createCollection() {
            const self = this;
            self.loading = true;

            const hashes = [];
            for( const hash of Object.keys(self.value)) {
                hashes.push(hash);
            }

            const requestParams = {
                method: 'POST',
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify({
                    Name: "Test " + new Date(),
                    Owner: "Tester",
                    Description: "Collection created at " + new Date(),
                    AssetHashes: hashes
                })
            }

            self.offset = 0;
            fetch('/collections/add', requestParams)
                .then(res => res.json())
                .then(json => {
                    console.log(json)
                });

        }
    },

    emits: ['componentEvent'],
}
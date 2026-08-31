export default {
    template: `
        <div>
            <div @dragover="dragOver" @drop="drop">
                <label class="btn btn-outline-secondary" role="button" for="formFileMultiple">Datei hochladen oder hierher ziehen</label>
                <input class="d-none" type="file" @change="fileChanged" id="formFileMultiple" multiple>
            </div>
            <div v-if="message" :class="'alert mt-2 ' + messageClass">
                {{ message }}
            </div>
            <div v-if="addedFiles && addedFiles.length" class="overflow-scroll" style="max-height: 60vh;">
                <div v-for="file in addedFiles" class="p-2"
                    role="button"
                    @click="assetClick(file)">
                    <template v-for="(origin, index) in file.Origins">
                        <div v-if="index == 0"><strong>{{ origin.Name }}</strong></div>
                        <div v-else class="small"> {{ origin.Name }}</div>
                    </template>
                </div>
            </div>
        </div>
    `,

    data() {
        return {
            message: '',
            messageClass: 'alert-success',

            addedFiles: [],
        }
    },

    methods: {

        fileChanged(e) {
            const self = this;
            const files = e.target.files || e.dataTransfer.files;
            if (!files.length)
                return;

            self.addedFiles.splice(0, self.addedFiles.length);
            self.message = "Start uploading";
            self.messageClass = 'alert-primary';

            for (const file of files) {

                self.message = "Uploading: " + file.name;
                self.messageClass = 'alert-primary';

                const request = new XMLHttpRequest();
                request.open('POST', '/assets/upload');
                request.upload.onprogress = function(progress) {
                    self.message = "Upload progress: " + progress.loaded + " of " + progress.total;
                    self.messageClass = 'alert-primary';
                }
                request.onreadystatechange = function () {
                    if (request.readyState === 4) {
                        if (request.status === 200) {

                            const json = JSON.parse(request.responseText)
                            console.log(json)

                            self.addUploadedFile(json, file)

                        } else if (this.status !== 0)
                            console.error("Upload failed", request)
                    }
                }
                request.send(file);
            }
        },
        addUploadedFile(json, file) {
            const self = this;

            self.message = "Adding file to archive: " + file.name;
            self.messageClass = 'alert-primary';

            const query = {
                TempName: json.tempName,
                Path: 'SPA/Upload',
                Name: file.name,
                Owner: "spa", //TODO
                FileTime: new Date(file.lastModified).toJSON()
            }

            const requestOptions = {
                method: 'POST',
                headers: {"Content-Type": "application/x-www-form-urlencoded"},
                body: JSON.stringify(query)
            }
            fetch('/assets/upload/add', requestOptions)
                .then(res => res.json())
                .then(json => {
                    for(const item of json) {
                        self.addedFiles.push(item);
                    }
                    self.message = "Added " + self.addedFiles.length + " file(s)";
                    self.messageClass = 'alert-success';
                });
        },
        dragOver(e) {
            e.preventDefault();
        },
        drop(e) {
            e.preventDefault();
            this.fileChanged(e);
        },
        assetClick(asset) {
            this.$emit('componentEvent', 'assetClick', 'upload', asset);
        }
    },

    emits: ['componentEvent'],
}
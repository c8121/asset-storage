(function () {
    return {

        template: `
            <div :class="'toast ' + showToastCss" role="alert" aria-live="assertive" aria-atomic="true">
                <div class="toast-header">
                    <strong class="me-auto">{{ headerCaption }}</strong>
                    <button type="button" class="btn-close" aria-label="Close" @click="hideToast"></button>
                </div>
                <div class="toast-body">
                    <div style="overflow: auto;max-height: 50vh;">
                        <div @dragover="dragOver" @drop="drop" :class="innerCssClass">
                            <label class="btn btn-secondary" role="button" for="formFileMultiple" v-html="labelCaption"></label>
                            <input class="d-none" type="file" @change="fileChanged" id="formFileMultiple" multiple>
                        </div>
                        <div v-if="message" :class="'alert mt-2 ' + messageClass">
                            {{ message }}
                        </div>
                        <div v-if="addedFiles && addedFiles.length">
                            <div v-for="file in addedFiles" class="p-2"
                                role="button"
                                @click="onMetaDataClick(file)">
                                <template v-for="(origin, index) in file.Origins">
                                    <div v-if="index == 0"><strong>{{ origin.Name }}</strong></div>
                                    <div v-else class="small"> {{ origin.Name }}</div>
                                </template>
                            </div>
                        </div>
                    </div>
                </div>
            </div>`,

        props: {
            headerCaption: {
                type: String,
                default: "Upload"
            },
            labelCaption: {
                type: String,
                default: "Datei hochladen oder hierher ziehen"
            },
            innerCssClass: {
                type: String,
                default: "p-3 border border-primary-subtle bg-body-tertiary"
            }
        },

        data() {
            return {
                message: '',
                messageClass: 'alert-success',

                addedFiles: [],

                showToastCss: ''
            }
        },
        methods: {
            showToast() {
                this.showToastCss = 'show';
                this.message = null;
                this.addedFiles.splice(0, this.addedFiles.length);
            },
            hideToast() {
                this.showToastCss = '';
            },
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
                    Name: file.name,
                    Owner: "spa", //TODO
                    FileTime: new Date(file.lastModified).toJSON()
                }
                client.post("/assets/upload/add", query).then((json) => {
                    for(const item of json) {
                        self.addedFiles.push(item);
                    }
                    self.message = "Added " + self.addedFiles.length + " file(s)";
                    self.messageClass = 'alert-success';
                })
            },
            dragOver(e) {
                e.preventDefault();
            },
            drop(e) {
                e.preventDefault();
                this.fileChanged(e);
            },
            onMetaDataClick(asset) {
                this.$emit('metaDataClick', asset);
            }
        },
        emits: [
            'uploadFinished',
            'metaDataClick'
        ],
    }
})();
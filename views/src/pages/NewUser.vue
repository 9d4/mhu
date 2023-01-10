<template>
    <div class="container" :class="{ 'd-none': loading }">
        <form @submit.prevent="onSubmit">
            <div class="mb-3">
                <p>Select user profile</p>
                <select class="form-select" v-model="selectedProfile">
                    <option v-for="profile in userprofiles">
                        {{ profile.name }}
                    </option>
                </select>
            </div>

            <div class="mb-3">
                <p>Select server</p>
                <select class="form-select" v-model="selectedServer">
                    <option v-for="server in servers">
                        {{ server.name }}
                    </option>
                </select>
            </div>

            <div class="mb-3">
                <div class="mb-3">
                    <label for="file" class="form-label">Upload users in csv</label>
                    <input class="form-control form-control-sm" id="file" type="file">
                </div>
            </div>

            <button>Push</button>
        </form>
    </div>
</template>

<script>
import agents from '../agents'

export default {
    data: () => ({
        loading: true,
        userprofiles: [],
        servers: [],

        selectedProfile: "",
        selectedServer: "",
    }),

    created() {
        agents.userprofiles().then(({ res, raw }) => {
            if (res === null) {
                alert("Error: ", raw.statusText)
                return
            }

            this.userprofiles = res;
        })

        agents.servers().then(({ res, raw }) => {
            if (res === null) {
                alert("Error: ", raw.statusText)
                return
            }

            this.servers = res;
        })

        this.loading = false;
    },

    methods: {
        onSubmit() {
            let file = document.getElementById('file').files[0]
            let xhr = new XMLHttpRequest
            let formData = new FormData

            formData.append('file', file)
            formData.append('server', this.selectedServer)
            formData.append('profile', this.selectedProfile)

            xhr.upload.addEventListener('progress', this.progressHandler, false)
            xhr.addEventListener('load', this.onLoadHandler, false)
            xhr.open('POST', '/api/upload', true)
            xhr.send(formData)
        },

        progressHandler(event) {
            //your code to track upload progress
            var p = Math.floor(event.loaded / event.total * 100);
            document.title = p + '%';
        },

        onLoadHandler(event) {
            // your code on finished upload
            alert(event.target.status);
        }
    }
}
</script>

<style>

</style>
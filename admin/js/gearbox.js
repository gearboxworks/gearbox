(function(){
    let app = new Vue({
        el: '#app',
        data: {
            apis: {},
            projects: [
                { "name": 'Some project 1' },
                { "name": 'Some project 2' }
            ],
        },
        methods: {
            getApiUrls() {
                this.$http.get('/api.json').then(function(response) {
                    this.apis = response.data
                });
            }
        }
    });
})();
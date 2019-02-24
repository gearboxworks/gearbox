function getGearboxProjects() {
    // if ( "undefined" == typeof gearbox ) {
    //     window.setTimeout(getGearboxProjects,50);
    //     return;
    // }
    //gearbox.loadProjects();
    // if (!gearbox.projects) {
    //     window.setTimeout(getGearboxProjects,50);
    //     return;
    // }
    var app = new Vue({
        el: '#app',
        data: {
            projects: [
                { "name": 'Some project 1' },
                { "name": 'Some project 2' }
            ],
        },
    });
    //alert('All good!');
}
getGearboxProjects();
